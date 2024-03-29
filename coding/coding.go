package coding

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"time"
)

func fixMistake(mistake byte, data byte) byte {
	switch mistake {
	case 1:
		return data ^ 1
	case 2:
		return data ^ 2
	case 4:
		return data ^ 4
	case 3:
		return data ^ 8
	case 6:
		return data ^ 16
	case 7:
		return data ^ 32
	case 5:
		return data ^ 64
	default:
		return data

	}
}
func rcrEncode(data byte) byte {
	num := data << 3
	buf := num
	pol := byte(11 << 3)
	var res byte
	for i := 0; i < 4; i++ {
		if bits.Len8(buf) < bits.Len8(pol) {
			//fmt.Println(bits.Len8(buf), bits.Len8(pol))
			pol = pol >> 1
			continue
		}
		res = buf ^ pol
		pol = pol >> 1
		buf = res
		//fmt.Printf("Result %d: %03b, generating polinom: %b \n", i, res, pol)
	}
	//fmt.Printf("coded string: %b", num|res)
	return num | res
}
func decode(data []byte) byte {

	//buf := data

	for i, buf := range data {
		pol := byte(11) << 3
		var res byte
		for i := 0; i < 4; i++ {
			if bits.Len8(buf) < bits.Len8(pol) && bits.Len8(buf) >= 4 {
				//fmt.Println(bits.Len8(buf), bits.Len8(pol))
				pol = pol >> 1
				continue
			}
			if bits.Len8(buf) < 4 {
				res = buf
				break
			}
			res = buf ^ pol
			pol = pol >> 1
			buf = res
			if bits.Len8(pol) < 4 {
				break
			}
			if res == 0 {
				break
			}
			//fmt.Printf("Result %d: %03b, generating polinom: %b \n", i, res, pol)

		}

		if res == 0 {
			data[i] = data[i] >> 3
		} else {
			data[i] = fixMistake(res, data[i]) >> 3
			fmt.Println("res: ", res, " buf: ", buf, " pol: ", pol)
			fmt.Printf("mistake syndrom :%b \n", res)
			fmt.Printf("%d fixed mistake:%b - %c \n", i, data[i], data[i])
		}

	}
	return data[0]<<4 | data[1]

}

// будем считать, что нам пришел байт с информацией.
// наш циклический код может закодировать только 4 бита, а закодированная инф-я будет занимать 7 бит
// предлагаю выделить первую и вторю четверку битов(обнулением всех остальных), прогнать через алгоритм циклического кодирования
// в итоге из исходного 1 байта получим 2 байта, в начале каждого будет 1 нулевой бит, который не используем.
func encode(data byte) (result []byte) {
	b1 := data & 240 >> 4
	b2 := data & 15
	//fmt.Println("halfs", b1, "  ", b2)
	return []byte{rcrEncode(b1), rcrEncode(b2)}

}
func makeMistake(data byte) (result byte) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(100) % 10

	if num == 0 {

		data = data ^ byte(math.Pow(2, float64(r.Intn(7))))
	}

	return data
}
func ProcessMessage(msg string) (res string, err error) {
	var processedMsg []byte
	var data []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Intn(100)%20 == 0 {
		return "", errors.New("lost message")
	}
	fmt.Println("MESSAGE ", msg)

	for i := 0; i < len(msg); i++ {
		data = encode(msg[i])
		fmt.Println("BYTE NUMBER ", i)
		fmt.Printf("initial data: %08b - %c \n", msg[i], msg[i])
		fmt.Printf("encoded data 1: %b \n", data[0])
		fmt.Printf("encoded data 2: %b \n", data[1])
		data[0] = makeMistake(data[0])
		data[1] = makeMistake(data[1])
		fmt.Printf("made mistake:%b,  %b \n", data[0], data[1])
		decodedData := decode(data)
		processedMsg = append(processedMsg, decodedData)
		fmt.Printf("decoded data: %08b - %c\n", decodedData, decodedData)
	}
	fmt.Println(processedMsg)
	fmt.Println(string(processedMsg))
	return string(processedMsg), nil
}
