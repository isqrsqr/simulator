package switch
import (
    "fmt"
    "net"
	"i2s2/media"
)



type Switch struct {
    Ports []*media.NetAccess
}

