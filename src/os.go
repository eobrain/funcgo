package os
import type (
	java.io.FileInputStream
)

func Open(path String){
	new FileInputStream(path)
}
