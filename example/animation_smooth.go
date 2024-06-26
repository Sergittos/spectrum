package main

import (
	"image/color"

	"github.com/cooldogedev/spectrum"
	"github.com/cooldogedev/spectrum/server"
	"github.com/cooldogedev/spectrum/session/animation"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
)

type smoothProcessor struct {
	consumer func(mgl32.Vec3, float32)
}

func (p *smoothProcessor) ProcessIncoming(packet.Packet) bool { return true }
func (p *smoothProcessor) ProcessOutgoing(pk packet.Packet) bool {
	if pk, ok := pk.(*packet.PlayerAuthInput); ok {
		p.consumer(pk.Position.Add(mgl32.Vec3{0, 1, 0}), pk.HeadYaw)
	}
	return true
}

func main() {
	logger := logrus.New()
	listenConfig := minecraft.ListenConfig{StatusProvider: spectrum.NewStatusProvider("Spectrum Proxy")}
	proxy := spectrum.NewSpectrum(server.NewStaticDiscovery(":19133"), logger, nil)
	if err := proxy.Listen(listenConfig); err != nil {
		logger.Errorf("Failed to listen on proxy: %v", err)
		return
	}

	for {
		s, err := proxy.Accept()
		if err != nil {
			logger.Errorf("Failed to accept session: %v", err)
			continue
		}

		s.SetAnimation(&animation.Smooth{
			Colour: color.RGBA{},
			Timing: protocol.CameraFadeTimeData{
				FadeInDuration:  0.75,
				WaitDuration:    3.25,
				FadeOutDuration: 0.75,
			},
		})
		s.SetProcessor(&smoothProcessor{consumer: func(position mgl32.Vec3, yaw float32) {
			s.Animation().(*animation.Smooth).Position = position
			s.Animation().(*animation.Smooth).Yaw = yaw
		}})
	}
}
