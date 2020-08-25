package utils

import (
	"bytes"
	"encoding/json"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

type IPFSConnect struct {
	Addr string
}

func (con *IPFSConnect) Add(content string) (string, error) {
	sh := shell.NewShell(con.Addr)
	cid, err := sh.Add(strings.NewReader(content))
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (con *IPFSConnect) BatchAdd(contents string) ([]string, error) {
	type Record struct {
		Key   string
		Value string
	}
	jsonStream := []byte(contents)
	var records []Record
	err := json.Unmarshal(jsonStream, &records)
	if err != nil {
		return nil, err
	}

	var cids []string
	sh := shell.NewShell(con.Addr)
	for _, record := range records {
		jsonStr, err := json.Marshal(record)
		if err != nil {
			return nil, err
		}
		cid, err := sh.Add(bytes.NewReader(jsonStr))
		if err != nil {
			return nil, err
		}
		cids = append(cids, cid)
	}
	return cids, nil
}

func (con *IPFSConnect) Get(cid string) (string, error) {
	sh := shell.NewShell(con.Addr)
	content, err := sh.Cat(cid)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	return buf.String(), nil
}

func (con *IPFSConnect) BatchGet(cids []string) ([]string, error) {
	var contents []string
	sh := shell.NewShell(con.Addr)
	for _, cid := range cids {
		content, err := sh.Cat(cid)
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(content)
		contents = append(contents, buf.String())
	}
	return contents, nil
}
