#include "mainwindow.h"

#include <QApplication>
#include<wiringPi.h>
#include<QFile>
#include<QIODevice>
#include<wiringPi.h>
#include"loginwindow.h"
//fan 29
 QString m_strProductKey;
 QString m_strDeviceName;
 QString m_strDeviceSecret;
 QString m_strRegionId;
 QString userId;
void init(){
    QFile file("config.json");
    if(!file.open(QFile::ReadOnly|QIODevice::Text)){
        qDebug()<<"fail";
    }
    QByteArray val=file.readAll();

  //  qDebug()<<val;
    QJsonParseError jsonError;
    QJsonDocument doucment = QJsonDocument::fromJson(val,&jsonError);  // 转化为 JSON 文档
    if(jsonError.error!=QJsonParseError::NoError){
        qDebug()<<"cao";
    }

    if (!doucment.isNull()&&jsonError.error==QJsonParseError::NoError) {  // 解析未发生错误

        if(doucment.isObject()){
            QJsonObject object=doucment.object();

            if(object.contains("DeviceName")){
                QJsonValue value=object.value("DeviceName");
                if(value.isString()){
                    m_strDeviceName=value.toString();
                  //  qDebug()<<m_strDeviceName;
                }
            }
            if(object.contains("ProductKey")){
                QJsonValue value=object.value("ProductKey");
                if(value.isString()){
                    m_strProductKey=value.toString();
                }
            }
            if(object.contains("DeviceSecret")){
                QJsonValue value=object.value("DeviceSecret");
                if(value.isString()){
                    m_strDeviceSecret=value.toString();
                }
            }
            if(object.contains("RegionId")){
                QJsonValue value=object.value("RegionId");
                if(value.isString()){
                    m_strRegionId=value.toString();
                }
            }
            if(object.contains("UserId")){
                QJsonValue value=object.value("UserId");
                if(value.isString()){
                    userId=value.toString();
                   // qDebug()<<userId;

                }
            }
        }
    }
    file.close();
}
int main(int argc, char *argv[])
{

    QApplication a(argc, argv);
    if(wiringPiSetup()==-1){
        qDebug()<<"setup wiringpi failed";
    }
//    if(wiringPiSetup()==-1){
//        qDebug()<<"setup wiringpi failed";
//    }

   init();

   qDebug()<<userId;
    if(userId==""){
        qDebug()<<userId;
        LoginWindow *w=new LoginWindow;
        w->show();

    }
    else{
        //qDebug()<<userId<<"!@#";
        MainWindow *w=new MainWindow;
        w->show();
    }

    return a.exec();
}


