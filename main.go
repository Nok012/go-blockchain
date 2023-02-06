package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Index	 int
	Hash     []byte
	Data     Tenant
	PrevHash []byte
}

type Tenant struct {
	Name      string
	Gender    string
	Email     string
	Phone     string
	RentHouse RentHouse
}

type RentHouse struct {
	HouseName    string
	Bedrooms     int64
	Bathrooms    int64
	Neighborhood string
	Submarket    string
	Borough      string
	Size         string
	Rent         int 
}

var rentHouse = []RentHouse{
	{
		HouseName:    "The Gables",
		Bedrooms:     1,
		Bathrooms:    1,
		Neighborhood: "Midtown",
		Submarket:    "All Midtown",
		Borough:      "Manhattan",
		Size:         "916 sqft",
		Rent:         4500,
	},
	{
		HouseName:    "Hillside",
		Bedrooms:     1,
		Bathrooms:    1,
		Neighborhood: "Astoria",
		Submarket:    "Northwest Queens",
		Borough:      "Queens",
		Size:         "996 sqft",
		Rent:         3500,
	},
	{
		HouseName:    "Sunnyside",
		Bedrooms:     1,
		Bathrooms:    2,
		Neighborhood: "Tribeca",
		Submarket:    "All Downtown",
		Borough:      "Manhattan",
		Size:         "1600 sqft",
		Rent:         10000,
	},
	{
		HouseName:    "Foxmoor Hall",
		Bedrooms:     2,
		Bathrooms:    2,
		Neighborhood: "Central Park South",
		Submarket:    "All Midtown",
		Borough:      "Manhattan",
		Size:         "1200 sqft",
		Rent:         5800,
	},
	{
		HouseName:    "Oaklands",
		Bedrooms:     1,
		Bathrooms:    1,
		Neighborhood: "Hamilton Heights",
		Submarket:    "All Upper Manhattan",
		Borough:      "Manhattan",
		Size:         "687 sqft",
		Rent:         2150,
	},
}

// convert struct to bytes array 
func EncodeToBytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// format json
func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "   "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func FormatData(data any) string{
	res1, _ := json.Marshal(data)
	res2, _ := PrettyString(string(res1))
	return res2
}



func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{[]byte(strconv.Itoa(b.Index)),EncodeToBytes(b.Data), b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(index int,data Tenant, prevHash []byte) *Block {
	block := &Block{ index, []byte{}, data, prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data Tenant) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	

	new := CreateBlock(len(chain.blocks)+1, data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

func Genesis() *Block {
	return CreateBlock(1, Tenant{}, []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func (chain *BlockChain) ChangeData(indexBlock int, fieldData interface{}, data string, rentHouse interface{}) {
	haveRentHouse := reflect.TypeOf(rentHouse) == reflect.TypeOf(RentHouse{})
	haveFieldData := reflect.TypeOf(fieldData) == reflect.TypeOf("")
	strFieldData := fmt.Sprintf("%v", fieldData)
	strFieldDataNoSpace := strings.TrimSpace(strFieldData)
	changeData := &chain.blocks[indexBlock-1].Data

	if  haveRentHouse && haveFieldData && len(strFieldDataNoSpace) != 0 {
		reflect.ValueOf(changeData).Elem().FieldByName(strFieldDataNoSpace).Set(reflect.ValueOf(data))
		changeData.RentHouse = rentHouse.(RentHouse)
	}else if haveRentHouse && !haveFieldData && len(strFieldDataNoSpace) != 0{
		changeData.RentHouse = rentHouse.(RentHouse)

	}else if !haveRentHouse && haveFieldData && len(strFieldDataNoSpace) != 0{
		reflect.ValueOf(changeData).Elem().FieldByName(strFieldDataNoSpace).Set(reflect.ValueOf(data))
	} 
}

func SetupBlockChain() *BlockChain {
	tenant := []Tenant{
		{	Name: "Helen", 	Gender: "Female", Email: "Helen123@gmail.com", 	Phone: "0842566597", RentHouse: rentHouse[0] },
		{	Name: "Elsie", 	Gender: "Female", Email: "Elsie123@gmail.com", 	Phone: "0652589651", RentHouse: rentHouse[1] },
		{	Name: "James", 	Gender: "Male",	  Email: "James123@gmail.com", 	Phone: "0826541857", RentHouse: rentHouse[2] },
		{	Name: "Steven", Gender: "Male",   Email: "Steven123@gmail.com", Phone: "0632585694", RentHouse: rentHouse[3] },
		{	Name: "Edward", Gender: "Male",   Email: "Edward123@gmail.com", Phone: "0611471417", RentHouse: rentHouse[4] },
		{	Name: "Linda", 	Gender: "Female", Email: "Linda123@gmail.com", 	Phone: "0815556966", RentHouse: rentHouse[0] },
		{	Name: "Sandra", Gender: "Fmale",  Email: "Sandra123@gmail.com", Phone: "0826541857", RentHouse: rentHouse[1] },
		{	Name: "Anna", 	Gender: "Fmale",  Email: "Anna123@gmail.com", 	Phone: "0965455532", RentHouse: rentHouse[2] },
		{	Name: "Peter", 	Gender: "Male",   Email: "Peter123@gmail.com", 	Phone: "0954771919", RentHouse: rentHouse[3] },
	}
	chain := InitBlockChain()
	chain.AddBlock(tenant[0])
	chain.AddBlock(tenant[1])
	chain.AddBlock(tenant[2])
	chain.AddBlock(tenant[3])
	chain.AddBlock(tenant[4])
	chain.AddBlock(tenant[5])
	chain.AddBlock(tenant[6])
	chain.AddBlock(tenant[7])
	chain.AddBlock(tenant[8])

	return chain
}



func (chain *BlockChain) ListBlock() *BlockChain {
	for i := 0; i <= len(chain.blocks)-1; i++ {
		chain.blocks[i].DeriveHash()
		fmt.Printf("Block: %d\n", chain.blocks[i].Index) 
		fmt.Printf("Previous Hash: %x\n", chain.blocks[i].PrevHash)
		fmt.Printf("Hash: %x\n", chain.blocks[i].Hash)
		fmt.Println("Data in Block: ", FormatData(chain.blocks[i].Data))
		fmt.Println()

		if i >= 2 {
			if string(chain.blocks[i].PrevHash) != string(chain.blocks[i-1].Hash) {
				fmt.Println("!!!!! Found data change !!!!!")
				fmt.Println("Data in Block:", i, "has been modified.")
				fmt.Println("")
				break
			}
		}
	}
	return chain
}

func (chain *BlockChain) GetWhoRentHouse(HouseName string) string {
	var name string
	for i := len(chain.blocks)-1; i >= 0; i-- {
		if chain.blocks[i].Data.RentHouse.HouseName == HouseName {
			fmt.Println("Latest tenant is name:",chain.blocks[i].Data.Name )
			name = chain.blocks[i].Data.Name
			break
		}
	}
	return name
}

func (chain *BlockChain) GetTenant(Name string) Tenant {	
	var tenant Tenant	
	for i := len(chain.blocks)-1; i >= 0; i-- {
		if chain.blocks[i].Data.Name == Name {
			fmt.Println("Tenant",FormatData(chain.blocks[i].Data) )
			tenant = chain.blocks[i].Data
			break
		}
	}
	return tenant
}

func (chain *BlockChain) ListWhoRentHouse(HouseName string) []Tenant {
	tenant := []Tenant{}
	for i := len(chain.blocks)-1; i >= 0; i-- {
		if chain.blocks[i].Data.RentHouse.HouseName == HouseName {
			tenant = append(tenant, chain.blocks[i].Data)
		}
	}
	ct := tenant[0]
	fmt.Printf("RentHouseName: %s\n",HouseName)
	fmt.Printf("Current tenant: %s\n",ct.Name)
	fmt.Printf("All tenanted: %v\n",FormatData(tenant))

	return tenant	
}

func (chain *BlockChain) ListAmountEverRentHouse(HouseName string) int {
	amount := 0
	for i := len(chain.blocks)-1; i >= 0; i-- {
		if chain.blocks[i].Data.RentHouse.HouseName == HouseName {
			amount++
		}
	}
	fmt.Printf("Amount have been rented: %d\n",amount)

	return amount
}



func main() {
	
	// Setup BlockChain ------------------------------------------------------
	chain := SetupBlockChain()

	// Add Block --------------------------------------------------------
	t := Tenant{	
		Name: "Nok", 
		Gender: "Male", 
		Email: "Nok012@gmail.com", 
		Phone: "0657019654", 
		RentHouse: rentHouse[0],	
	}
	chain.AddBlock(t)

	// Edit Data) ------------------------------------------------------------
	chain.ChangeData(2,"Name","Nok",rentHouse[4])
	// List BlockChain ---------------------------------------------------------	
	chain.ListBlock()


	chain.ListWhoRentHouse("Sunnyside")
	chain.GetWhoRentHouse("Sunnyside")
	chain.ListAmountEverRentHouse("Sunnyside")
	chain.GetTenant("Anna")
	
}
	

	
	

	
	
	


	


