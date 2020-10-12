package common

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/spf13/viper"
	"pack.ag/amqp"
	"server/models"
	"time"
)

var (
	uid             string
	accessKey       string
	accessSecret    string
	region          string
	consumerGroupId string
	clientId        string
	endpoint        string
)

var i = 1

func InitAmqp() {
	uid = viper.GetString("aliyuniot.uid")
	accessKey = viper.GetString("aliyuniot.accessKey")
	accessSecret = viper.GetString("aliyuniot.accessSecret")
	region = viper.GetString("aliyuniot.region")
	consumerGroupId = viper.GetString("aliyuniot.consumerGroupId")
	clientId = viper.GetString("aliyuniot.clientId")
	endpoint = "https://" + uid + ".iot-as-http2." + region + ".aliyuncs.com"

	address := fmt.Sprintf("amqps://%s.iot-amqp.%s.aliyuncs.com:5671", uid, region)
	timestamp := time.Now().Nanosecond() / 1000000
	username := fmt.Sprintf("%s|authMode=aksign,signMethod=Hmacsha1,consumerGroupId=%s,authId=%s,timestamp=%d|",
		clientId, consumerGroupId, accessKey, timestamp)
	stringToSign := fmt.Sprintf("authId=%s&timestamp=%d", accessKey, timestamp)
	hmacKey := hmac.New(sha1.New, []byte(accessSecret))
	hmacKey.Write([]byte(stringToSign))
	password := base64.StdEncoding.EncodeToString(hmacKey.Sum(nil))

	amqpManager := &AmqpManager{
		address:  address,
		username: username,
		password: password,
	}

	//如果需要做接受消息通信或者取消操作 从Background衍生context
	ctx := context.Background()

	amqpManager.startReceiveMessage(ctx)
}

//业务函数。用户自定义实现，该函数被异步执行，请考虑系统资源消耗情况。
func (am *AmqpManager) processMessage(message *amqp.Message) {
	//fmt.Println("data received:", string(message.GetData()), " properties:", message.ApplicationProperties)
	//msg:=string(message.GetData())
	//fmt.Println(msg)

	deviceFilter(string(message.GetData()))

}

type AmqpManager struct {
	address  string
	username string
	password string
	client   *amqp.Client
	session  *amqp.Session
	receiver *amqp.Receiver
}

func (am *AmqpManager) startReceiveMessage(ctx context.Context) {

	childCtx, _ := context.WithCancel(ctx)
	err := am.generateReceiverWithRetry(childCtx)
	//print(err)
	if nil != err {
		return
	}
	defer func() {
		am.receiver.Close(childCtx)
		am.session.Close(childCtx)
		am.client.Close()
	}()

	for {

		//阻塞接受消息，如果ctx是background则不会被打断。
		message, err := am.receiver.Receive(ctx)

		if nil == err {
			go am.processMessage(message)
			message.Accept()
		} else {
			fmt.Println("amqp receive data error:", err)

			//如果是主动取消，则退出程序。
			select {
			case <-childCtx.Done():
				return
			default:
			}

			//非主动取消，则重新建连。
			err := am.generateReceiverWithRetry(childCtx)
			//	print(err)
			if nil != err {
				return
			}

		}
	}

}

func (am *AmqpManager) generateReceiverWithRetry(ctx context.Context) error {

	//退避重试 10ms 依次x2 直到 20s
	duration := 10 * time.Millisecond
	maxDuration := 20000 * time.Millisecond
	times := 1

	//异常情况，退避重连。
	for {
		select {
		case <-ctx.Done():
			return amqp.ErrConnClosed
		default:
		}

		err := am.generateReceiver()
		if nil != err {
			time.Sleep(duration)
			if duration < maxDuration {
				duration *= 2
			}
			fmt.Println("amqp connect retry,times:", times, ",duration:", duration)
			times++
		} else {
			fmt.Println("amqp connect init success")
			return nil
		}
	}
}

//由于包不可见，无法判断conn和session状态，重启连接获取。
func (am *AmqpManager) generateReceiver() error {

	//topic和credit在此处没有实际作用，云端没有定义这两个参数。
	if am.session != nil {
		receiver, err := am.session.NewReceiver(
			amqp.LinkSourceAddress("/queue-name"),
			amqp.LinkCredit(20),
		)
		//如果断网等行为发生，conn会关闭导致session建立失败，未关闭连接则建立成功。
		if err == nil {
			am.receiver = receiver
			return nil
		}
	}

	//清理上一个连接。
	if am.client != nil {
		am.client.Close()
	}

	client, err := amqp.Dial(am.address, amqp.ConnSASLPlain(am.username, am.password))
	if err != nil {
		return err
	}
	am.client = client

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	am.session = session

	receiver, err := am.session.NewReceiver(
		amqp.LinkSourceAddress("/queue-name"),
		amqp.LinkCredit(20),
	)
	if err != nil {
		return err
	}
	am.receiver = receiver

	return nil
}

func deviceFilter(message string) {
	//db:=GetDB()
	fmt.Println(message)
	j := gjson.New(message)
	productKey := j.GetString("productKey")
	iotId := j.GetString("iotId")
	gmtCreate := j.GetString("gmtCreate")
	//items := j.GetString("items")
	//status := j.GetString("status")
	//	deviceName := j.GetString("deviceName")
	temperature := j.GetFloat32("items.temperature.value")
	humidity := j.GetFloat32("items.humidity.value")
	lightIntensity := j.GetInt("items.light_intensity.value")
	soil := j.GetInt("items.soil.value")
	//	fan_switch:=j.GetBool("items.fan_switch.value")
	//light_switch:=j.GetBool("items.light_switch.value")
	//water_switch:=j.GetBool("items.water_switch.value")

	switch productKey {
	//pi
	case "a1I7rPDpEx5":
		piHistoryData := models.PiHistoryData{
			IotId:          iotId,
			Time:           gmtCreate,
			Temperature:    temperature,
			Humidity:       humidity,
			LightIntensity: lightIntensity,
			SoilMoisture:   soil,
		}
		year := time.Now().Year()
		month := time.Now().Month()
		fmt.Println(piHistoryData)
		piname := fmt.Sprintf("%v-%s-pi", year, month)
		fmt.Println(piname)
	//piname:= strconv.Itoa(year)+"-"+ strconv.Itoa(month)+"-pi"
	//err:=db.Create(&piHistoryData).Table(piname)

	default:

	}
}
