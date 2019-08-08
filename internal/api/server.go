package server

import (
	"encoding/json"
	"fmt"
	"net"

	google_protobuf1 "github.com/golang/protobuf/ptypes/empty"
	"github.com/pythonrocks/ustd-example/internal/service"
	context "golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type rpcServer struct {
	env *service.Env
	rpc *service.USTDCLient
}

// StartRPCServer runs RPC server
func StartRPCServer(env *service.Env) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.Conf.RPCPort))
	if err != nil {
		return err
	}

	s := rpcServer{env, service.NewClient(env)}

	grpcServer := grpc.NewServer()

	RegisterExampleAPIServer(grpcServer, s)

	grpcServer.Serve(lis)
	return nil
}

// ListAddresses returns list of addresses for user's wallet
func (s rpcServer) ListAddresses(in *google_protobuf1.Empty, stream ExampleAPI_ListAddressesServer) error {
	result, err := s.rpc.Call("getaddressesbyaccount", "")
	if err != nil {
		return grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}

	addresses := []string{}
	if err := json.Unmarshal(result.Result, &addresses); err != nil {
		return grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}

	for _, addr := range addresses {
		result, err = s.rpc.Call("omni_getallbalancesforaddress", addr)
		if err != nil {
			return grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
		}
		aInfo := []struct {
			PropertyID int32  `json:"propertyid"`
			Name       string `json:"name"`
			Balance    string `json:"balance"`
			Reserved   string `json:"reserved"`
			Frozen     string `json:"frozen"`
		}{}

		if err := json.Unmarshal(result.Result, &aInfo); err != nil {
			return grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
		}

		balances := []*Balance{}
		for _, b := range aInfo {
			balances = append(balances, &Balance{b.PropertyID, b.Name, b.Balance, b.Reserved, b.Frozen})
		}
		if err := stream.Send(&AddressInfo{addr, balances}); err != nil {
			return grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
		}
	}
	return nil
}

// GetTransactionInfo returns info about transaction.
func (s rpcServer) GetTransactionInfo(ctx context.Context, in *TransactionInfoRequest) (*TransactionInfo, error) {
	result, err := s.rpc.Call("omni_gettransaction", in.TxID)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}
	tx := struct {
		Txid             string `json:"txid"`
		Sendingaddress   string `json:"sendingaddress"`
		Referenceaddress string `json:"referenceaddress"`
		Ismine           bool   `json:"ismine"`
		Confirmations    int32  `json:"confirmations"`
		Fee              string `json:"fee"`
		Blocktime        int32  `json:"blocktime"`
		Valid            bool   `json:"valid"`
		Positionblock    int32  `json:"positionblock"`
		Version          int32  `json:"version"`
		TypeInt          int32  `json:"type_int"`
		Type             string `json:"type"`
	}{}
	if err := json.Unmarshal(result.Result, &tx); err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}
	return &TransactionInfo{
		tx.Txid,
		tx.Sendingaddress,
		tx.Referenceaddress,
		tx.Ismine,
		tx.Confirmations,
		tx.Fee,
		tx.Blocktime,
		tx.Valid,
		tx.Positionblock,
		tx.Version,
		tx.TypeInt,
		tx.Type,
	}, nil
}

// GetWalletInfo returns info with user's wallet balances.
func (s rpcServer) GetWalletInfo(ctx context.Context, in *google_protobuf1.Empty) (*WalletInfo, error) {
	result, err := s.rpc.Call("omni_getwalletbalances")
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}
	wInfo := []struct {
		PropertyID int32  `json:"propertyid"`
		Name       string `json:"name"`
		Balance    string `json:"balance"`
		Reserved   string `json:"reserved"`
		Frozen     string `json:"frozen"`
	}{}
	if err := json.Unmarshal(result.Result, &wInfo); err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}
	balances := []*Balance{}
	for _, b := range wInfo {
		balances = append(balances, &Balance{b.PropertyID, b.Name, b.Balance, b.Reserved, b.Frozen})
	}
	return &WalletInfo{
		Balances: balances,
	}, nil
}

// GetAddressInfo can be used to obtain address's balances
func (s rpcServer) GetAddressInfo(ctx context.Context, in *AddressInfoRequest) (*AddressInfo, error) {

	result, err := s.rpc.Call("omni_getallbalancesforaddress", in.Address)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}

	aInfo := []struct {
		PropertyID int32  `json:"propertyid"`
		Name       string `json:"name"`
		Balance    string `json:"balance"`
		Reserved   string `json:"reserved"`
		Frozen     string `json:"frozen"`
	}{}
	if err := json.Unmarshal(result.Result, &aInfo); err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}
	balances := []*Balance{}
	for _, b := range aInfo {
		balances = append(balances, &Balance{b.PropertyID, b.Name, b.Balance, b.Reserved, b.Frozen})
	}

	return &AddressInfo{in.Address, balances}, nil
}

// NewAddress allows to generate new address
func (s rpcServer) NewAddress(ctx context.Context, in *google_protobuf1.Empty) (*AddressInfo, error) {
	result, err := s.rpc.Call("getnewaddress")
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}

	address := ""
	if err := json.Unmarshal(result.Result, &address); err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}

	result, err = s.rpc.Call("omni_getallbalancesforaddress", address)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}

	aInfo := []struct {
		PropertyID int32  `json:"propertyid"`
		Name       string `json:"name"`
		Balance    string `json:"balance"`
		Reserved   string `json:"reserved"`
		Frozen     string `json:"frozen"`
	}{}
	if err := json.Unmarshal(result.Result, &aInfo); err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}

	balances := []*Balance{}
	for _, b := range aInfo {
		balances = append(balances, &Balance{b.PropertyID, b.Name, b.Balance, b.Reserved, b.Frozen})
	}

	return &AddressInfo{address, balances}, nil
}

// Transfer lets user to send currency to another address.
func (s rpcServer) Transfer(ctx context.Context, in *TransferRequest) (*TransferResult, error) {
	result, err := s.rpc.Call("omni_send", in.Fromaddress, in.Toaddress, in.Propertyid, in.Amount)
	if err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}

	hash := ""
	if err := json.Unmarshal(result.Result, &hash); err != nil {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("%s", err))
	}

	return &TransferResult{hash}, nil
}
