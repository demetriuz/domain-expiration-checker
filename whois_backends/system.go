package whois_backends

import "os/exec"

type SystemWhoisBackend struct{}


func (b SystemWhoisBackend) Fetch(domain string) (string, error){
	out, err := exec.Command("whois", domain).Output()
	if err != nil{
		return "", err
	}else{
		return string(out[:]), nil
	}
}
