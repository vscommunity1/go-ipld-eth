package ipldeth

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	//common "github.com/ethereum/go-ethereum/common"

	rlp "github.com/ethereum/go-ethereum/rlp"
)

func TestTrieParsing(t *testing.T) {
	v := "f874822080b86ff86d078504a817c80083015f90943d4a3fdbb4ffae950a069ac7319f157bbaaa010e8810a741a462780000801ca05b7175be69fccc145074194c187e4ad3e13f2b50a13042efa7b55f8346305f05a02d008fff5b7169f4aaa85a405359e4ce023cd9fd52c87a0492a1359aaebfce68"

	data, err := hex.DecodeString(v)
	if err != nil {
		t.Fatal(err)
	}

	tn, err := NewTrieNode(data)
	if err != nil {
		t.Fatal(err)
	}

	_ = tn
}

func TestWeirdCase(t *testing.T) {
	v := "e218a0a396a1abf4d8dc066e724a3b00a9650d49c1b37674875ce0e25c2c8df2e75674"
	data, err := hex.DecodeString(v)
	if err != nil {
		t.Fatal(err)
	}

	tn, err := NewTrieNode(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tn.Links())
}

func TestHash(t *testing.T) {
	fi, err := os.Open("test_data/eth-state-trie-root-4052365-0xa94ca535ef7f51d70728c15c7a292dc2fec7915949117166558eb012e2dda92d")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(fi)

	tn, err := NewTrieNode(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	cid := tn.Cid()
	if cid.String() != "z45oqTS5yKSHRBwfqEPFFd63Dugf6QxreYc8h7MF3CpxCpFhkFv" {
		t.Fatal("Wrong calculated cid: %v", cid)
	}

	hash := tn.(*TrieNode).Hash()
	if hash != "a94ca535ef7f51d70728c15c7a292dc2fec7915949117166558eb012e2dda92d" {
		t.Fatal("Expected hashes to match")
	}

}

func incrParse(data []byte) {
	k, val, rest, err := rlp.Split(data)
	if err != nil {
		panic(err)
	}
	if k == rlp.List {
		fmt.Println("[")
		incrParse(val)
		fmt.Println("]")
	} else if k == rlp.String {
		if len(val) > 3 {
			tx, err := ParseTx(val)
			if err != nil {
				fmt.Println("wasnt a tx:", err)
			} else {
				fmt.Println(tx)
			}
		} else {
			fmt.Println("keyval:", val)
		}
	} else {
		fmt.Println("thing:", k, val)
	}
	if len(rest) > 0 {
		fmt.Println("rest:", rest, len(rest))
		incrParse(rest)
	}
}
