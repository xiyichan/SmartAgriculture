#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include<QString>
#include<QtMqtt/qmqttclient.h>
#include<QJsonDocument>
#include<QJsonObject>
#include<QJsonValue>
#include<QJsonParseError>
#include<QMessageBox>
#include<QMessageAuthenticationCode>
#include"qrencode.h"
#include<QFile>
#include<QIODevice>
#include<wiringPi.h>
#include<QDebug>
#include<QString>
#include<QTimer>
#include<QNetworkAccessManager>
#include<QNetworkRequest>
#include<QNetworkReply>
#define HIGH_TIME 32

QT_BEGIN_NAMESPACE
namespace Ui { class MainWindow; }
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(QWidget *parent = nullptr);
    ~MainWindow();
private slots:
    void on_pushButton_water_clicked();

    void on_pushButton_fan_clicked();

    void on_pushButton_light_clicked();

    void finishedSlot(QNetworkReply *reply);
    void on_pushButton_loginout_clicked();



    void on_pushButton_auto_save_clicked();

    void on_pushButton_2_clicked();

private:
    Ui::MainWindow *ui;
    int dht11=7;
    int water=0;
    int fan=29;
    int light=25;
    int lightv=24;
    unsigned long data_dht11;
    QString temperature;
    QString humidity;
    bool readDht11Data();
    double ads1115_voltage[4];
    int16_t ads1115_value[4];
    void timerEvent(QTimerEvent *event);

    bool water_switch=false;
    bool fan_switch=false;
    bool light_switch=false;
    bool auto_switch=false;
    int auto_water=20000;
    int auto_light=20000;
    int auto_fan=20;
    QMqttClient *m_client;
//        QString m_strProductKey;
//        QString m_strDeviceName;
//        QString m_strDeviceSecret;
//        QString m_strRegionId;

    QString m_strPubTopic;
    QString m_strSubTopic;
    QString m_strTargetServer;
    QString payload2;
    int id=0;

    QJsonParseError simp_json_error;
    QJsonDocument simp_parse_doucement;
    QString method;
    void parse(QString message);
    void setFanSwitch(bool s);
    void setWaterSwitch(bool s);
    void setLightSwitch(bool s);
    void setAutoSwitch(bool s);
    QNetworkAccessManager *m_accessManager;

};
#endif // MAINWINDOW_H
