package gtoken

import (
	"fmt"
	"gitee.com/johng/gf/g/encoding/gjson"
	jwt "github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

var (
	key []byte = []byte("Hello World！This is secret!")
)

// type Token struct {
// 	Raw       string                 // The raw token.  Populated when you Parse a token
// 	Method    SigningMethod          // The signing method used or to be used
// 	Header    map[string]interface{} // The first segment of the token
// 	Claims    Claims                 // The second segment of the token
// 	Signature string                 // The third segment of the token.  Populated when you Parse a token
// 	Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
// }

// type StandardClaims struct {
// 	Audience  string `json:"aud,omitempty"`
// 	ExpiresAt int64  `json:"exp,omitempty"`
// 	Id        string `json:"jti,omitempty"`
// 	IssuedAt  int64  `json:"iat,omitempty"`
// 	Issuer    string `json:"iss,omitempty"`
// 	NotBefore int64  `json:"nbf,omitempty"`
// 	Subject   string `json:"sub,omitempty"`
// }

// 产生json web token
func GenToken(username interface{}) string {
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix()),
		ExpiresAt: int64(time.Now().Unix() + 1000),
		Issuer:    "admin",
		Subject:   username.(string),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println("生成的TOKEN", ss)
	return ss
}

//&Token{
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
// eyJleHAiOjE1MzE5MTAzNTMsImlzcyI6ImFkbWluIiwibmJmIjoxNTMxOTA5MzUzLCJzdWIiOiJ5eXgxIn0.
// D5vxn3ffUOaE330BPUAwIpxrnHxg0dKlSyzwpHyXxLg

//  0xc0421b7420

//  map[alg:HS256 typ:JWT]

//   map[exp:1.531910353e+09 iss:admin nbf:1.531909353e+09 sub:yyx1]

//   D5vxn3ffUOaE330BPUAwIpxrnHxg0dKlSyzwpHyXxLg  	true}
//验证
func CheckToken(tokenstr string) (bool, string) {
	token, err := jwt.Parse(tokenstr, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return false, err.Error()
	} else {
		raw := strings.Split(token.Raw, ".")
		claims, _ := jwt.DecodeSegment(raw[1]) //Claims部分
		claims_json, _ := gjson.DecodeToJson(claims)
		return true, claims_json.Get("sub").(string)
	}
	//token.raw.claims => {"exp":1531913061,"iss":"admin","nbf":1531912061,"sub":"yyx1"}

}
