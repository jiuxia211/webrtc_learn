package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jiuxia211/webrtc_learn/signal"
	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/webrtc/v4"
)

// 简单的示例
func main() {
	// 创建一个空的音视频编解码器
	m := &webrtc.MediaEngine{}
	// 注册编解码器(codecs)
	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		// VP8编码器，时钟速率，SDP格式参数行，RTCP反馈
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8, ClockRate: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		//  RTP 数据包的有效载荷类型，视频编码器
		PayloadType: 96}, webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}

	// 注册默认拦截器
	i := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		panic(err)
	}
	// PLI 拦截器
	intervalPliFactory, err := intervalpli.NewReceiverInterceptor()
	if err != nil {
		panic(err)
	}
	i.Add(intervalPliFactory)
	// 根据 编解码器和拦截器创建api
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))

	// 创建一个ICE服务器
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				// 这是 Google 提供的一个 STUN（Session Traversal Utilities for NAT）服务器，用于进行 NAT 穿透，帮助设备在 NAT 网络环境中建立对等连接。
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	// 创建peerConnection
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if cErr := peerConnection.Close(); cErr != nil {
			fmt.Printf("cannot close peerConnection: %v\n", cErr)
		}
	}()
	// 创建一个我们输出视频到浏览器的通道
	outputTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	if err != nil {
		panic(err)
	}
	// 把Track加入peerConnection
	rtpSender, err := peerConnection.AddTrack(outputTrack)
	if err != nil {
		panic(err)
	}
	// 开启协程 不断读取 RTCP数据包
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()
	// 读取来自浏览器的SessionDescription 并解析(TODO)
	offer := webrtc.SessionDescription{}
	signal.Decode(signal.MustReadStdin(), &offer)

	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// 最重要的一个handler，控制视频轨道的读取和发送
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, r *webrtc.RTPReceiver) {
		fmt.Printf("Track has started, of type %d: %s \n", track.PayloadType(), track.Codec().MimeType)

		for {
			rtp, _, readErr := track.ReadRTP()
			if readErr != nil {
				panic(readErr)
			}

			if writeErr := outputTrack.WriteRTP(rtp); writeErr != nil {
				panic(writeErr)
			}
		}
	})
	// 控制数据的传输
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())
			for range time.NewTicker(5 * time.Second).C {
				message := "test send message"
				fmt.Printf("Sending '%s'\n", message)
				// Send the message as text
				sendErr := d.SendText(message)
				if sendErr != nil {
					panic(sendErr)
				}
			}
		})
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})
	// 监测Peer Connection的状态
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("Peer Connection State has changed: %s\n", s.String())
		if s == webrtc.PeerConnectionStateFailed {
			fmt.Println("Peer Connection has gone to failed exiting")
			os.Exit(0)
		}

		if s == webrtc.PeerConnectionStateClosed {
			fmt.Println("Peer Connection has gone to closed exiting")
			os.Exit(0)
		}
	})
	// 创建一个answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// 阻塞直到ICE complete完成
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))

	// Block forever
	select {}

}
