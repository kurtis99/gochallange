package drum

import "bytes"
import "encoding/binary"
import "io/ioutil"
import "fmt"

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	p := &Pattern{}

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(path)
	}

	buf := bytes.NewReader(dat)

	if err := binary.Read(buf, binary.LittleEndian, &p.magic); err != nil {
		panic("Failed to read binary data: " + err.Error())
	}

	if err := binary.Read(buf, binary.BigEndian, &p.len); err != nil {
		panic("Failed to read binary data: " + err.Error())
	}

	if err := binary.Read(buf, binary.LittleEndian, &p.hw_version); err != nil {
		panic("Failed to read binary data: " + err.Error())
	}

	if err := binary.Read(buf, binary.LittleEndian, &p.tempo); err != nil {
		panic("Failed to read binary data: " + err.Error())
	}

	fmt.Printf("%s\n", string(p.magic[:]))
	fmt.Printf("%d\n", p.len)
	fmt.Printf("%s\n", string(p.hw_version[:]))
	fmt.Printf("%f\n", p.tempo)

	for buf.Len() > 0 {
		t := &Track{}

		if err := binary.Read(buf, binary.LittleEndian, &t.index); err != nil {
			panic("Failed to read binary data: " + err.Error())
		}

		if err := binary.Read(buf, binary.LittleEndian, &t.name_len); err != nil {
			panic("Failed to read binary data: " + err.Error())
		}

		t.name = make([]byte, t.name_len)

		if err := binary.Read(buf, binary.LittleEndian, &t.name); err != nil {
			panic("Failed to read binary data: " + err.Error())
		}

		if err := binary.Read(buf, binary.LittleEndian, &t.drums.d); err != nil {
			panic("Failed to read binary data: " + err.Error())
		}

		fmt.Print(t)
	}

	return p, nil
}

func (t *Track) String() string {
	return fmt.Sprintf("(%d) %s \t%s\n", t.index, t.name, t.drums.String())
}

func (t *Pattern) String() string {
	return fmt.Sprintf("%s\n", "huita")
}

func (d *Drums) String() string {

	out := make([]byte, 21)

	for i, v := range d.d {

		if i%4 == 0 {
			out = append(out, '|')
		}

		if v != 0 {
			out = append(out, 'x')
		} else {
			out = append(out, '-')
		}
	}

	out = append(out, '|')

	return string(out)
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	magic      [6]byte
	len        uint64
	hw_version [32]byte
	tempo      float32
	track      []Track
}

type Track struct {
	index    uint32
	name_len uint8
	name     []byte
	drums    Drums
}

type Drums struct {
	d [16]byte
}
