package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

// RPCRequest represent a RCP request
type RPCRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int64       `json:"id"`
	JSONRpc string      `json:"jsonrpc"`
}

// RPCError represents an error
type RPCError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Guarantee RPCError satisfies the builtin error interface.
var _, _ error = RPCError{}, (*RPCError)(nil)

// Error returns a string describing the RPC error.
func (e RPCError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// RPCResponse represents response object.
type RPCResponse struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    *RPCError       `json:"error"`
}

// USTDCLient represents json-rpc client
type USTDCLient struct {
	env        *Env
	client     *http.Client
	serverAddr string
}

func paramsWrap(params ...interface{}) interface{} {
	var finalParams interface{}

	if params != nil {
		switch len(params) {
		case 0:
		case 1:
			if params[0] != nil {
				var typeOf reflect.Type

				for typeOf = reflect.TypeOf(params[0]); typeOf != nil && typeOf.Kind() == reflect.Ptr; typeOf = typeOf.Elem() {
				}

				if typeOf != nil {
					switch typeOf.Kind() {
					case reflect.Struct, reflect.Array, reflect.Slice, reflect.Interface, reflect.Map:
						finalParams = params[0]
					default:
						finalParams = params
					}
				}
			} else {
				finalParams = params
			}
		default:
			finalParams = params
		}
	}

	return finalParams
}

// NewClient returns rpc client
func NewClient(env *Env) *USTDCLient {
	return &USTDCLient{
		env:        env,
		serverAddr: fmt.Sprintf("http://%s:%s", env.Conf.Host, env.Conf.Port),
		client:     &http.Client{},
	}
}

// Call does json-rpc call and returns result.
func (c *USTDCLient) Call(method string, params ...interface{}) (RPCResponse, error) {
	var result RPCResponse
	request := RPCRequest{
		Method:  method,
		Params:  paramsWrap(params...),
		ID:      time.Now().UnixNano(),
		JSONRpc: "1.0"}
	payloadBuffer := &bytes.Buffer{}
	jsonEncoder := json.NewEncoder(payloadBuffer)
	err := jsonEncoder.Encode(request)
	if err != nil {
		return result, err
	}
	req, err := http.NewRequest("POST", c.serverAddr, payloadBuffer)
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")

	if len(c.env.Conf.RPCUser) > 0 || len(c.env.Conf.RPCPassword) > 0 {
		req.SetBasicAuth(c.env.Conf.RPCUser, c.env.Conf.RPCPassword)
	}

	log.Printf("RPC request %#v\n", req)
	resp, err := c.client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)
	log.Printf("RPC response %#v\n", result)
	return result, nil
}
