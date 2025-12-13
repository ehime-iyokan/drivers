//go:build atsamd51 || atsame5x

package ili9341

import (
	"device/sam"
	"errors"
	"machine"
	"unsafe"

	tinygoDma "github.com/sago35/tinygo-dma"
)

type spiDriver struct {
	bus *machine.SPI
}

func NewSPI(bus *machine.SPI, dc, cs, rst machine.Pin) *Device {
	return &Device{
		dc:  dc,
		cs:  cs,
		rst: rst,
		rd:  machine.NoPin,
		driver: &spiDriver{
			bus: bus,
		},
	}
}

// NewSPI() を元にして、DMAの設定を追加したもの
func NewSPIWithDMA(bus *machine.SPI, dc, cs, rst machine.Pin, dmaSrcBuf *[]byte) (*Device, error) {
	dma := tinygoDma.NewDMA(func(d *tinygoDma.DMA) {
		return
	})
	dma.SetTrigger(tinygoDma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM7_TX)
	dma.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	buf := *dmaSrcBuf
	if len(buf) == 0 {
		return &Device{}, errors.New("dma src buf is empty")
	}
	desc := dma.GetDescriptor()
	desc.UpdateDescriptor(tinygoDma.DescriptorConfig{
		SRC:    unsafe.Pointer(&buf[0]),
		DST:    unsafe.Pointer(&bus.Bus.DATA.Reg),
		SRCINC: tinygoDma.DMAC_SRAM_BTCTRL_SRCINC_ENABLE,
		DSTINC: tinygoDma.DMAC_SRAM_BTCTRL_DSTINC_DISABLE,
		SIZE:   uint32(len(buf)),
	})

	return &Device{
		dc:  dc,
		cs:  cs,
		rst: rst,
		rd:  machine.NoPin,
		driver: &spiDriver{
			bus: bus,
		},
		dma: dma,
	}, nil
}

func (pd *spiDriver) configure(config *Config) {
}

func (pd *spiDriver) write8(b byte) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPIM_CTRLB_RXEN)

	for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
	}
	pd.bus.Bus.DATA.Set(uint32(b))

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPIM_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write8n(b byte, n int) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPIM_CTRLB_RXEN)

	for i, c := 0, n; i < c; i++ {
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(b))
	}

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPIM_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write8sl(b []byte) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPIM_CTRLB_RXEN)

	for i, c := 0, len(b); i < c; i++ {
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(b[i]))
	}

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPIM_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write16(data uint16) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPIM_CTRLB_RXEN)

	for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
	}
	pd.bus.Bus.DATA.Set(uint32(uint8(data >> 8)))
	for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
	}
	pd.bus.Bus.DATA.Set(uint32(uint8(data)))

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPIM_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write16n(data uint16, n int) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPIM_CTRLB_RXEN)

	for i := 0; i < n; i++ {
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data >> 8)))
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data)))
	}

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPIM_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write16sl(data []uint16) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPIM_CTRLB_RXEN)

	for i, c := 0, len(data); i < c; i++ {
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data[i] >> 8)))
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data[i])))
	}

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPIM_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_CTRLB) {
	}
}
