#include "loginwindow.h"
#include "ui_loginwindow.h"
#include"mainwindow.h"
extern QString m_strProductKey;
extern QString m_strDeviceName;
extern QString m_strDeviceSecret;
extern QString m_strRegionId;
extern QString userId;

LoginWindow::LoginWindow(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::LoginWindow)
{
    ui->setupUi(this);

    qDebug()<<m_strDeviceName;
    //genQrcodeImage("asdDD",160,160);
    //QPixmap pixmap=QPixmap::fromImage(*image);
    //ui->label->setPixmap(pixmap);
    ui->label->setPixmap(generateQR(m_strDeviceName,320,320));
    m_client=new QMqttClient(this);
    QString m_strSubTopic = "/"+m_strProductKey + "/" + m_strDeviceName + "/user/get";//订阅topic
    QString m_strTargetServer = m_strProductKey + ".iot-as-mqtt." + m_strRegionId + ".aliyuncs.com";//域名


    m_client->setHostname(m_strTargetServer);
    m_client->setPort(1883);
    QString clientId="d5R46fOSfNNwVTNRuSaM";         //表示客户端ID，建议使用设备的MAC地址或SN码，64字符内。
    QString signmethod = "hmacsha1";    //加密方式
    QString message ="clientId"+clientId+"deviceName"+m_strDeviceName+"productKey"+m_strProductKey;
    //  qDebug()<<message;
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
        QStringList list=topic.name().split("/");
        QString topic11111="/sys/"+m_strProductKey+"/"+m_strDeviceName+"/rrpc/response/"+list[list.size()-1];
        qDebug()<<topic11111;
        m_client->publish(topic11111,message);
        if(message!=""){
            QFile file("config.json");
            if(!file.open(QIODevice::ReadWrite)){
                qDebug()<<"111";
            }
            QString uid=message;
            QJsonObject jsonObject;
            jsonObject.insert("ProductKey",m_strProductKey);
            jsonObject.insert("DeviceName",m_strDeviceName);
            jsonObject.insert("DeviceSecret",m_strDeviceSecret);
            jsonObject.insert("RegionId",m_strRegionId);
            jsonObject.insert("UserId",uid);

            QJsonDocument jsonDoc;
            jsonDoc.setObject(jsonObject);
            file.write(jsonDoc.toJson());
            file.close();
            MainWindow *w=new MainWindow;
            w->show();
            this->close();
        }

    });
    connect(m_client, &QMqttClient::pingResponseReceived, this, [this]() {
        const QString content = QDateTime::currentDateTime().toString()
                + QLatin1String(" PingResponse")
                + QLatin1Char('\n');
        qDebug()<<content;
    });
    qDebug()<<m_strSubTopic;
    auto subscription = m_client->subscribe(m_strSubTopic);
    if (subscription) {
        QMessageBox::critical(this, QLatin1String("Error"), QLatin1String("Could not subscribe. Is there a valid connection?"));

    }
}
LoginWindow::~LoginWindow()
{
    delete ui;
}
QPixmap LoginWindow::createQRCode(const QString &text)
{
    int margin = 2;
    if (text.length() == 0)
    {
        return QPixmap();
    }
    QRcode *qrcode = QRcode_encodeString(text.toLocal8Bit(), 2, QR_ECLEVEL_L, QR_MODE_8, 0);
    if (qrcode == NULL) {
        return QPixmap();
    }
    unsigned char *p, *q;
    p = NULL;
    q = NULL;
    int x, y, bit;
    int realwidth;

    realwidth = qrcode->width;
    QImage image = QImage(realwidth, realwidth, QImage::Format_Indexed8);
    QRgb value;
    value = qRgb(255, 255, 255);
    image.setColor(0, value);
    value = qRgb(0, 0, 0);
    image.setColor(1, value);
    image.setColor(2, value);
    image.fill(0);
    p = qrcode->data;
    for (y = 0; y<qrcode->width; y++) {
        bit = 7;
        q += margin / 8;
        bit = 7 - (margin % 8);
        for (x = 0; x<qrcode->width; x++) {
            if ((*p & 1) << bit)
                image.setPixel(x, y, 1);
            else
                image.setPixel(x, y, 0);
            bit--;
            if (bit < 0)
            {
                q++;
                bit = 7;
            }
            p++;
        }
    }
    return QPixmap::fromImage(image.scaledToWidth(200));

}
QPixmap LoginWindow::generateQR(QString strContent ,int width,int height)
{
    //[1]
    QPixmap GenQRPixmap;
    QRcode *qrcode; //二维码数据
    qrcode  = QRcode_encodeString(strContent.toStdString().c_str(),2,QR_ECLEVEL_Q,QR_MODE_8,1);

    qint32 qrcode_width = qrcode->width > 0 ? qrcode->width : 1;

    //[2]
    double scale_x = (double)width / (double)qrcode_width; //二维码图片的缩放比例
    double scale_y =(double) height /(double) qrcode_width;
    QImage qrImg = QImage(width, height, QImage::Format_ARGB32);

    QPainter painter1(&qrImg);
    QColor background(Qt::white);
    painter1.setBrush(background);
    painter1.setPen(Qt::NoPen);
    painter1.drawRect(0, 0, width, height);
    QColor foreground(Qt::black);
    painter1.setBrush(foreground);

    for( qint32 y = 0; y < qrcode_width; y ++)
    {
        for(qint32 x = 0; x < qrcode_width; x++)
        {
            unsigned char b = qrcode->data[y * qrcode_width + x];
            if(b & 0x01)
            {
                QRectF r(x * scale_x, y * scale_y, scale_x, scale_y);
                painter1.drawRects(&r, 1);
            }
        }
    }

    //[3]
    GenQRPixmap = QPixmap::fromImage(qrImg);

    return GenQRPixmap;
}

