package drum

import "bytes"
import "encoding/binary"
import "io/ioutil"
import "fmt"

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {

	var diff int

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

	diff = buf.Len() - int(p.len)

	if err := binary.Read(buf, binary.LittleEndian, &p.hwVersion); err != nil {
		panic("Failed to read binary data: " + err.Error())
	}

	if err := binary.Read(buf, binary.LittleEndian, &p.tempo); err != nil {
		panic("Failed to read binary data: " + err.Error())
	}

	for buf.Len() > diff {
		t := &Track{}

		if err := binary.Read(buf, binary.LittleEndian, &t.index); err != nil {
			panic("Failed to read binary data: " + err.Error())
		}

		if err := binary.Read(buf, binary.LittleEndian, &t.nameLen); err != nil {
			panic("Failed to read binary data: " + err.Error())
		}

		t.name = make([]byte, t.nameLen)

		if err := binary.Read(buf, binary.LittleEndian, &t.name); err != nil {
			panic("Failed to read binary data: " + err.Error())
		}

		if err := binary.Read(buf, binary.LittleEndian, &t.drums.d); err != nil {
			panic("Failed to read binary data: " + err.Error())
		}

		p.track = append(p.track, t)
	}

	return p, nil
}

func (t *Track) String() string {

	return fmt.Sprintf("(%d) %s\t%s\n",
		t.index,
		t.name,
		t.drums.String())
}

func (p *Pattern) String() string {

	var b bytes.Buffer

	n := bytes.IndexByte(p.hwVersion[:], 0)
	fmt.Fprintf(&b, "Saved with HW Version: %s\n", string(p.hwVersion[:n]))
	fmt.Fprintf(&b, "Tempo: %g\n", p.tempo)

	for _, v := range p.track {
		fmt.Fprintf(&b, "%s", v.String())
	}

	return b.String()
}

func (d *Drums) String() string {

	var b [21]byte
	var offset int

	for i, v := range d.d {

		if i%4 == 0 {
			b[i+offset] = '|'
			offset++
		}

		if v != 0 {
			b[i+offset] = 'x'
		} else {
			b[i+offset] = '-'
		}
	}

	b[20] = '|'

	return fmt.Sprintf("%s", b)
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	magic     [6]byte
	len       uint64
	hwVersion [32]byte
	tempo     float32
	track     []*Track
}

// Track is contains information about single track
type Track struct {
	index   uint32
	nameLen uint8
	name    []byte
	drums   Drums
}

// Drums is helper struct to print drum pattern
type Drums struct {
	d [16]byte
}
