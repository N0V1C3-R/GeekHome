package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/dgrijalva/jwt-go"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
)

type Response struct {
	StatusCode   int
	ResponseText string
}

func GenerateSnowflake() int64 {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	unixTime := ConvertToNanoTime(GetCurrentTime()) + rand.Int63n(999)
	snowflakeId := node.Generate().Int64() + unixTime
	return snowflakeId
}

func GetToken(secretKey string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": ConvertToMilliTime(GetCurrentTime()),
		},
	)
	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}

func VerifyToken(tokenString, secretKey string) (bool, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return false, err
	}
	if !token.Valid {
		return false, nil
	}
	return true, nil
}

func SerializationObj(obj interface{}) string {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

func CreateFolder(filePath string) {
	if err := os.MkdirAll(filePath, 0777); err != nil {
		panic(err)
	}
}

func Base64EncodeString(stdIn []byte) string {
	return base64.StdEncoding.EncodeToString(stdIn)
}

func Base64DecodeString(stdIn string) (stdout []byte) {
	stdout, err := base64.StdEncoding.DecodeString(stdIn)
	if err != nil {
		return
	}
	return stdout
}

func IsBase64String(input string) bool {
	if len(input)%4 != 0 {
		return false
	}
	base64Pattern := "^[A-Za-z0-9+/]*={0,2}$"
	match := CheckStringWithRegex(input, base64Pattern)
	return match
}

func IsNumeric(input string) bool {
	regExp := regexp.MustCompile(`^[0-9]+$`)

	return regExp.MatchString(input)
}

func ConvertWithSign(number string, fromBase, toBase int) (string, error) {
	n := new(big.Int)
	_, success := n.SetString(number, fromBase)
	if !success {
		return "", fmt.Errorf("invalid number in base %d: %s", fromBase, number)
	}

	sign := ""
	if n.Sign() == -1 {
		sign = "-"
		n.Abs(n)
	}

	result := ""
	for n.Sign() > 0 {
		remainder := new(big.Int)
		n.QuoRem(n, big.NewInt(int64(toBase)), remainder)
		result = string(rune('0'+int(remainder.Int64()))) + result
	}

	if result == "" {
		result = "0"
	}

	return sign + result, nil
}

func GenerateVerificationCode() string {
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}

func LimitDecimalPlaces(num float64, decimalPlaces int) float64 {
	format := fmt.Sprintf("%%.%df", decimalPlaces)
	result, _ := strconv.ParseFloat(fmt.Sprintf(format, num), 64)
	return result
}

func CheckStringWithRegex(str, allowedChars string) bool {
	if str == "" {
		return false
	}
	match, _ := regexp.MatchString(allowedChars, str)
	return match
}

func RemoveElements(slice, elem interface{}) interface{} {
	sliceValue := reflect.ValueOf(slice)
	elemValue := reflect.ValueOf(elem)

	result := reflect.MakeSlice(sliceValue.Type(), 0, sliceValue.Len())

	for i := 0; i < sliceValue.Len(); i++ {
		e := sliceValue.Index(i)
		if !reflect.DeepEqual(e.Interface(), elemValue.Interface()) {
			result = reflect.Append(result, e)
		}
	}

	return result.Interface()
}

func CallFunction(funcName interface{}, args ...interface{}) interface{} {
	fnValue := reflect.ValueOf(funcName)
	argValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		argValues[i] = reflect.ValueOf(arg)
	}
	resultValues := fnValue.Call(argValues)
	if len(resultValues) > 0 {
		return resultValues[0].Interface()
	}
	return nil
}
