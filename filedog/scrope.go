package filedog

import (
	"fmt"
)
// this is used to test globle vars in the same package.if here is defined,then in the other files of the same package,
// no more same vars should be defined.the compilor will check it
var globleHosts string

func Show() {
	fmt.Println(globleHosts)
}
