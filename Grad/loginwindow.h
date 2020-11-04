#ifndef LOGINWINDOW_H
#define LOGINWINDOW_H

#include <QWidget>
#include<QJsonDocument>
#include<QJsonObject>
#include<QImage>
#include<QPainter>
#include<QJsonValue>
#include<QJsonParseError>
#include<QMessageBox>
#include<QMessageAuthenticationCode>
#include"qrencode.h"
#include<QFile>
#include<QIODevice>
#include<QDebug>
#include<QNetworkReply>
#include<QtMqtt/qmqttclient.h>

namespace Ui {
class LoginWindow;
}

class LoginWindow : public QWidget
{
    Q_OBJECT

public:
    explicit LoginWindow(QWidget *parent = nullptr);
    ~LoginWindow();
    QPixmap createQRCode(const QString &text);
    QPixmap generateQR(QString strContent ,int width,int height);

private slots:


private:
QMqttClient *m_client;
    Ui::LoginWindow *ui;
    QNetworkAccessManager *m_accessManager;

};

#endif // LOGINWINDOW_H
