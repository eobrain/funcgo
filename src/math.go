package math
import(
)

Pi := Math::PI

func Nextafter(x double, y double) {
	java.lang.Math::nextAfter(x, y)
}
