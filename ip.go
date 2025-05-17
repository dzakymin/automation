package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

var (
	octet3        *string
	gateway       *string
	interfaceName *string
)

type InputOctet func() (*string, *string)

func flagiOctet() (*string, *string) {

	regex := regexp.MustCompile(`^(?:[1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-4])(?:\.(?:[1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-4])){3}/(24|25|26|27)$`)
	regex2 := regexp.MustCompile(`^(?:[1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-4])(?:\.(?:[1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-4])){3}$`)
	if regex.MatchString(*octet3) && regex2.MatchString(*gateway) {
		return octet3, gateway
	} else {
		return nil, nil
	}

}

func editIP(ipEdit InputOctet, interfaceName string) string {
	flag.Parse()
	ip, gateway := ipEdit()
	if ip == nil && gateway == nil {
		return "address failed"
	} else {
		edit_ip := fmt.Sprintf(`network:
  version: 2
  ethernets:
    %s:
      addresses:
      - %s
      nameservers:
        addresses:
        - 8.8.8.8
        search: []
      routes:
      - to: "default"
        via: %s`, interfaceName, *ip, *gateway)

		return edit_ip
	}
}

func main() {
	octet3 = flag.String("ip", "172.28.1.1/25", "Insert ip of this machine")
	gateway = flag.String("gateway", "172.28.10.254", "Insert gateway")
	interfaceName = flag.String("interface", "ens18", "Interface name")

	flag.Parse()
	path := "/etc/netplan/50-cloud-init.yaml"
	write_file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("error pak le", err.Error())
	}
	defer write_file.Close()

	result := editIP(flagiOctet, *interfaceName)

	if result != "address failed" {
		write_file.WriteString(result)
		fmt.Println("Success write to file")
		cmd := exec.Command("sudo", "netplan", "apply")
		cmd.Stdout = os.Stdout
		cmd.Run()
		if err != nil {
			fmt.Println("Error while applying command")
		} else {
			fmt.Println("Network configuration up")
		}
	} else {
		fmt.Println("error masbro")
	}

}
