#include "mainwindow.h"
#include "ui_mainwindow.h"
#include "loginwindow.h"
#include"ads1115.h"
extern QString m_strProductKey;
extern QString m_strDeviceName;
extern QString m_strDeviceSecret;
extern QString m_strRegionId;
extern QString userId;
MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
{

    ui->setupUi(this);
    m_accessManager=new QNetworkAccessManager(this);
    //qDebug()<<"123123"<<m_strDeviceName;
    ui->pushButton_water->setStyleSheet("background:rgb(0,255,0)");
    ui->pushButton_fan->setStyleSheet("background:rgb(0,255,0)");
    ui->pushButton_light->setStyleSheet("background:rgb(0,255,0)");
    ui->pushButton_2->setStyleSheet("background:rgb(0,255,0)");
    ui->spinBox_fan->setValue(auto_fan);
    ui->spinBox_water->setValue(auto_water);
    ui->spinBox_light->setValue(auto_light);
//    if(wiringPiSetup()==-1){
//        qDebug()<<"setup wiringpi failed";
//    }
    ads1115Setup(100,0x48);
    startTimer(5000);
    m_client=new QMqttClient(this);

//    m_strProductKey="a1I7rPDpEx5";  //需要跟阿里云Iot平台一致;
//    m_strDeviceName="d5R46fOSfNNwVTNRuSaM";   //需要跟阿里云Iot平台一致;
//    m_strDeviceSecret="1c7985af68584e861990d66cce938a02";   //需要跟阿里云平台一致
//    m_strRegionId="cn-shanghai";

    m_strPubTopic = "/sys/" + m_strProductKey + "/" + m_strDeviceName + "/thing/event/property/post";//发布topic
    m_strSubTopic = "/sys/" + m_strProductKey + "/" + m_strDeviceName + "/thing/service/property/set";//订阅topic
    m_strTargetServer = m_strProductKey + ".iot-as-mqtt." + m_strRegionId + ".aliyuncs.com";//域名
   QString testtopic="/"+m_strProductKey + "/" + m_strDeviceName + "/user/bind";//订阅topic
    m_client->setHostname(m_strTargetServer);
    m_client->setPort(1883);
    QString clientId="d5R46fOSfNNwVTNRuSaM";         //表示客户端ID，建议使用设备的MAC地址或SN码，64字符内。
    QString signmethod = "hmacsha1";    //加密方式
    QString message ="clientId"+clientId+"deviceName"+m_strDeviceName+"productKey"+m_strProductKey;

    m_client->setUsername(m_strDeviceName + "&" + m_strProductKey);

    m_client->setClientId(clientId + "|securemode=3,signmethod=" + signmethod + "|");
    // m_client->setPassword("36CF2FDDD251957716F2D6CCD52EB0741FF00FB4");
    m_client->setPassword(QMessageAuthenticationCode::hash(message.toLocal8Bit(),
                                                           m_strDeviceSecret.toLocal8Bit(),
                                                           QCryptographicHash::Sha1).toHex());

    m_client->connectToHost();




    connect(m_client, &QMqttClient::messageReceived, this, [this](const QByteArray &message, const QMqttTopicName &topic) {
        const QString content = QDateTime::currentDateTime().toString()
                + QLatin1String(" Received Topic: ")
                + topic.name()
                + QLatin1String(" Message: ")
                + message
                + QLatin1Char('\n');
        qDebug()<<content;
        //parseCH(message);
        parse(message);
    });
    //类似心跳包？一分钟一次
        connect(m_client, &QMqttClient::pingResponseReceived, this, [this]() {
            const QString content = QDateTime::currentDateTime().toString()
                    + QLatin1String(" PingResponse")
                    + QLatin1Char('\n');
            qDebug()<<content;
        });
        // payload2="{\"id\":1,\"params\": {\"CH1\":10},\"method\": \"thing.event.property.post\"}";

    //订阅
    qDebug()<<m_strSubTopic;
    //m_client->subscribe(testtopic);
    auto subscription = m_client->subscribe(m_strSubTopic);
    if (subscription) {
        QMessageBox::critical(this, QLatin1String("Error"), QLatin1String("Could not subscribe. Is there a valid connection?"));

    }
    qDebug()<<testtopic;
    //m_client->subscribe(testtopic);
    auto subscription1 = m_client->subscribe(testtopic);
    if (subscription1) {
        QMessageBox::critical(this, QLatin1String("Error"), QLatin1String("Could not subscribe. Is there a valid connection?"));

    }


    qDebug()<<QString::number(m_client->state());
ui->pushButton_light->setStyleSheet("background:rgb(0,255,0)");

    setWaterSwitch(0);
    setFanSwitch(0);
    setLightSwitch(0);

}
MainWindow::~MainWindow()
{

    delete ui;
}

bool MainWindow::readDht11Data(){
    unsigned char crc;
    unsigned char i;
    data_dht11=0;
    pinMode(dht11,OUTPUT);
    digitalWrite(dht11,0);
    delay(25);
    digitalWrite(dht11,1);
    pinMode(dht11,INPUT);
    pullUpDnControl(dht11, PUD_UP);

    delayMicroseconds(27);
    if (digitalRead(dht11) == 0) //SENSOR ANS
    {
        while (!digitalRead(dht11))
            ; //wait to high

        for (i = 0; i < 32; i++)
        {
            while (digitalRead(dht11))
                ; //data clock start
            while (!digitalRead(dht11))
                ; //data start
            delayMicroseconds(HIGH_TIME);
            data_dht11 *= 2;
            if (digitalRead(dht11) == 1) //1
            {
                data_dht11++;
            }
        }

        for (i = 0; i < 8; i++)
        {
            while (digitalRead(dht11))
                ; //data clock start
            while (!digitalRead(dht11))
                ; //data start
            delayMicroseconds(HIGH_TIME);
            crc *= 2;
            if (digitalRead(dht11) == 1) //1
            {
                crc++;
            }
        }
        return 1;
    }
    else
    {
        return 0;
    }
}

void MainWindow::timerEvent(QTimerEvent *event){
    pinMode(dht11,OUTPUT);
    digitalWrite(dht11,1);
    if(readDht11Data())
    {
        pinMode(dht11,OUTPUT);
        digitalWrite(dht11,1);
        temperature=QString::number((data_dht11>>8)&0xff,10)+"."+QString::number((data_dht11)&0xff,10);
        humidity=QString::number((data_dht11>>24)&0xff,10)+"."+QString::number((data_dht11>>16)&0xff,10);
        ui->label_humidity->setText(humidity);
        ui->label_temperature->setText(temperature);
    }
    else{

    }
    data_dht11=0;
    //ads1115
    for(int i=100;i<102;i++){
        ads1115_value[i-100]=(int16_t)analogRead(i);
        qDebug()<<"asd1115::"<<i-100<<ads1115_value[i-100];
        ads1115_voltage[i-100]=ads1115_value[i-100]*(4.096/32768);
        qDebug()<<"asd1115::"<<i-100<<ads1115_voltage[i-100];
    }
    ui->label_soil->setText(QString::number(ads1115_value[0],10));
    ui->label_lightIntensity->setText(QString::number(ads1115_value[1],10));
    if (auto_switch==true){
    if(temperature.toFloat()>auto_fan){
        setFanSwitch(1);
    }else{
        setFanSwitch(0);
    }
    if(ads1115_value[1]>auto_light){
          setLightSwitch(1);
    }else{
        setLightSwitch(0);
    }
    if(ads1115_value[0]>auto_water){
       setWaterSwitch(1);
    }else{
        setWaterSwitch(0);
    }
}
    QString payload=QString("{\"id\":1,\"params\": {\"temperature\":%1,\"humidity\":%2,\"light_intensity\":%3,\"soil\":%4,\"water_switch\""
                            ":%5,\"light_switch\":%6,\"fan_switch\":%7,\"auto_switch\":%8},\"method\": \"thing.event.property.post\"}").arg(temperature).arg(humidity)
            .arg(ads1115_value[1]).arg(ads1115_value[0]).arg(water_switch).arg(light_switch).arg(fan_switch).arg(auto_switch);
    m_client->publish(m_strPubTopic,payload.toUtf8());
}





void MainWindow::on_pushButton_water_clicked()
{
    if(water_switch==false){
        setWaterSwitch(1);
    }else{
        setWaterSwitch(0);
    }
}

void MainWindow::on_pushButton_fan_clicked()
{
    if(fan_switch==false){
        setFanSwitch(1);
    }else{
        setFanSwitch(0);
    }
}

void MainWindow::on_pushButton_light_clicked()
{
    if(light_switch==false){
        setLightSwitch(1);
    }else{
        setLightSwitch(0);
    }
}
void MainWindow::parse(QString message){
    QJsonParseError jsonError;
    QJsonDocument doucment = QJsonDocument::fromJson(message.toUtf8(),&jsonError);  // 转化为 JSON 文档
    if (!doucment.isNull()&&jsonError.error==QJsonParseError::NoError) {  // 解析未发生错误
        if(doucment.isObject()){
            QJsonObject object=doucment.object();
            if(object.contains("params")){
                QJsonValue value=object.value("params");
                if(value.isObject()){
                    QJsonObject obj=value.toObject();
                    if(obj.contains("water_switch")){
                        QJsonValue water=obj.value("water_switch");
                        if(water.isDouble()){
                            // water_switch=water.toDouble();
                            //qDebug()<<water.toDouble();
                            setWaterSwitch(water.toDouble());

                        }
                    }
                    if(obj.contains("fan_switch")){
                        QJsonValue fan=obj.value("fan_switch");
                        if(fan.isDouble()){
                            //fan_switch=fan.toDouble();
                            // qDebug()<<fan.toDouble();
                            setFanSwitch(fan.toDouble());
                        }
                    }
                    if(obj.contains("light_switch")){
                        QJsonValue l=obj.value("light_switch");
                        if(l.isDouble()){
                            //light_switch=l.toDouble();
                            // qDebug()<<l.toDouble();
                            setLightSwitch(l.toDouble());
                        }
                    }
                    if(obj.contains("auto_switch")){
                        QJsonValue l=obj.value("auto_switch");
                        if(l.isDouble()){
                            //light_switch=l.toDouble();
                            // qDebug()<<l.toDouble();
                            setAutoSwitch(l.toDouble());
                        }
                    }
                }
            }
        }

    }
}


void MainWindow::setFanSwitch(bool s){
    if(s==true){
        pinMode(fan,OUTPUT);
        digitalWrite(fan,1);
        ui->pushButton_fan->setText("fan_open");
        ui->pushButton_fan->setStyleSheet("background:rgb(255,0,0)");
        fan_switch=true;
    }
    else{
        pinMode(fan,OUTPUT);
        ui->pushButton_fan->setText("fan_close");
        ui->pushButton_fan->setStyleSheet("background:rgb(0,255,0)");
        digitalWrite(fan,0);
        fan_switch=false;
    }
}
void MainWindow::setLightSwitch(bool s){
    if(s==true){
        pinMode(lightv,OUTPUT);
        digitalWrite(lightv,1);

        ui->pushButton_light->setText("light_open");
        ui->pushButton_light->setStyleSheet("background:rgb(255,0,0)");
        light_switch=true;

    }
    else{
        pinMode(lightv,OUTPUT);
        digitalWrite(lightv,0);
        ui->pushButton_light->setText("light_close");
        ui->pushButton_light->setStyleSheet("background:rgb(0,255,0)");
        light_switch=false;
    }
}
void MainWindow::setWaterSwitch(bool s){
    if(s==true){
        pinMode(water,OUTPUT);
        digitalWrite(water,1);
        ui->pushButton_water->setText("water_open");
        ui->pushButton_water->setStyleSheet("background:rgb(255,0,0)");
        water_switch=true;
    }
    else{
        pinMode(water,OUTPUT);
        ui->pushButton_water->setText("water_close");
        ui->pushButton_water->setStyleSheet("background:rgb(0,255,0)");
        digitalWrite(water,0);
        water_switch=false;
    }
}





void MainWindow::on_pushButton_loginout_clicked()
{
    QFile file("config.json");
    if(!file.open(QIODevice::ReadWrite)){
        qDebug()<<"111";
    }

    QJsonObject jsonObject;
    jsonObject.insert("ProductKey",m_strProductKey);
    jsonObject.insert("DeviceName",m_strDeviceName);
    jsonObject.insert("DeviceSecret",m_strDeviceSecret);
    jsonObject.insert("RegionId",m_strRegionId);
    jsonObject.insert("UserId","");

    QJsonDocument jsonDoc;
    jsonDoc.setObject(jsonObject);
    file.write(jsonDoc.toJson());
    file.close();

 QNetworkRequest *request = new QNetworkRequest();
    //需要配置文件
    request->setUrl(QUrl("http://47.100.108.193:8080/api/device/pi/login/out"));

    QByteArray postData;
    QByteArray responseData;
    //表单


    QString deviceName="DeviceName="+m_strDeviceName+"&";
    QString deviceSecret="DeviceSecret="+m_strDeviceSecret+"&";

    postData.append(deviceName);
    postData.append(deviceSecret);
    //  postData.append("Email=2863768433@qq.com&Password=czx987852");

    QNetworkReply* reply = m_accessManager->post(*request,postData);

    //同步
    QEventLoop eventLoop;
    connect(m_accessManager, SIGNAL(finished(QNetworkReply*)), &eventLoop, SLOT(quit()));
    eventLoop.exec();       //block until finish
    responseData = reply->readAll();
    qDebug()<<responseData;

    QJsonParseError json_error;
    QJsonDocument parse_doucment = QJsonDocument::fromJson(responseData, &json_error);
qDebug()<<"114";
    if(json_error.error == QJsonParseError::NoError)

    {

        if(parse_doucment.isObject())
        {

            QJsonObject obj = parse_doucment.object();
                QJsonValue code_value = obj.take("code");
                int code = code_value.toInt();
                QJsonValue msg_value = obj.take("msg");
                QString msg=msg_value.toString();
                qDebug()<<msg;
            if(code==200){
               this->close();
            }
            else{
                QMessageBox::information(this,"提示",msg);
            }

        }

    }


//    LoginWindow *w=new LoginWindow;
//    w->show();
    //TODO:
    //you bug  huan hui login windows hui chu xian
    //wiringPiNewNode: Pin 100 overlaps with existing definition
    //only rset
   // this->close();
}
void MainWindow::finishedSlot(QNetworkReply *reply)
{
    if(reply->error()==QNetworkReply::NoError){
        QByteArray bytes=reply->readAll();
        qDebug()<<bytes;
    }else{
        qDebug("code : %d",(int)reply->error());
        qDebug(qPrintable(reply->errorString()));
    }
    reply->deleteLater();
}

void MainWindow::setAutoSwitch(bool s){
    if(s==true){
        pinMode(fan,OUTPUT);
        digitalWrite(fan,1);
        ui->pushButton_2->setText("atuo_open");
        ui->pushButton_2->setStyleSheet("background:rgb(255,0,0)");
        auto_switch=true;
    }
    else{
        pinMode(fan,OUTPUT);
        ui->pushButton_2->setText("atuo_close");
        ui->pushButton_2->setStyleSheet("background:rgb(0,255,0)");
        digitalWrite(fan,0);
        auto_switch=false;
        setWaterSwitch(0);
        setLightSwitch(0);
        setFanSwitch(0);
    }
}

void MainWindow::on_pushButton_auto_save_clicked()
{
    auto_water=ui->spinBox_water->value();
    auto_fan=ui->spinBox_fan->value();
    auto_light=ui->spinBox_light->value();
    qDebug()<<auto_water<<auto_fan<<auto_light;
}

void MainWindow::on_pushButton_2_clicked()
{
    if(auto_switch==false){
        setAutoSwitch(1);
    }else{
        setAutoSwitch(0);
    }
}
