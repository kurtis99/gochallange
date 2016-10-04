package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/nacl/box"
)

type SecureReader struct {
	r          io.Reader
	buf        []byte
	publicKey  *[32]byte
	privateKey *[32]byte
}

type SecureWriter struct {
	w          io.Writer
	buf        []byte
	publicKey  *[32]byte
	privateKey *[32]byte
}

func (sr *SecureReader) Read(p []byte) (n int, err error) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(sr.r)

	var nonce [24]byte
	copy(nonce[:], buf.Bytes()[:24])

	out, ok := box.Open(sr.buf, buf.Bytes()[24:], &nonce, sr.publicKey, sr.privateKey)
	if !ok {
		return 0, errors.New("SecureReader: failed to read (nonce huita)")
	}

	copy(p, out)

	return len(out), nil
}

func (sw *SecureWriter) Write(p []byte) (n int, err error) {

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic("Failed to get nonce: " + err.Error())
		return 0, errors.New("SecureWriter: nonce failed")
	}

	encrypted := box.Seal(nonce[:], p, &nonce, sw.publicKey, sw.privateKey)

	sw.w.Write(encrypted)

	return len(encrypted), nil
}

// NewSecureReader instantiates a new SecureReader
func NewSecureReader(r io.Reader, priv, pub *[32]byte) io.Reader {

	sr := &SecureReader{r: r, publicKey: pub, privateKey: priv}

	return sr
}

// NewSecureWriter instantiates a new SecureWriter
func NewSecureWriter(w io.Writer, priv, pub *[32]byte) io.Writer {

	sw := &SecureWriter{w: w, publicKey: pub, privateKey: priv}

	return sw
}

// Dial generates a private/public key pair,
// connects to the server, perform the handshake
// and return a reader/writer.
func Dial(addr string) (io.ReadWriteCloser, error) {
	return nil, nil
}

// Serve starts a secure echo server on the given listener.
func Serve(l net.Listener) error {
	return nil
}

func main() {
	port := flag.Int("l", 0, "Listen mode. Specify port")
	flag.Parse()

	// Server mode
	if *port != 0 {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
		log.Fatal(Serve(l))
	}

	// Client mode
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <port> <message>", os.Args[0])
	}
	conn, err := Dial("localhost:" + os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if _, err := conn.Write([]byte(os.Args[2])); err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, len(os.Args[2]))
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf[:n])
}
