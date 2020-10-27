#ifndef LOGINWINDOW_H
#define LOGINWINDOW_H

#include <QWidget>
#include<QJsonDocument>
#include<QJsonObject>
#include<QJsonValue>
#include<QJsonParseError>
#include<QMessageBox>
#include<QMessageAuthenticationCode>
#include"qrencode.h"
#include<QFile>
#include<QIODevice>
#include<QDebug>
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
private:

    Ui::LoginWindow *ui;
};

#endif // LOGINWINDOW_H
