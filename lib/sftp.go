package lib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Config represents SSH connection parameters.
type SftpConfig struct {
	Username     string
	Password     string
	PrivateKey   string
	Server       string
	KeyExchanges []string
	KeyPsw       string
	Timeout      time.Duration
}

// Client provides basic functionality to interact with a SFTP server.
type Client struct {
	config     SftpConfig
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

// New initialises SSH and SFTP clients and returns Client type to use.
func NewSftpclient(config SftpConfig) (*Client, error) {
	c := &Client{
		config: config,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return c, nil
}

// Create creates a remote/destination file for I/O.
func (c *Client) Create(filePath string) (io.ReadWriteCloser, error) {
	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return c.sftpClient.Create(filePath)
}
func (c *Client) OpenFile(filePath string) (io.ReadWriteCloser, error) {
	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return c.sftpClient.OpenFile(filePath, os.O_WRONLY)

}

// Upload writes local/source file data streams to remote/destination file.
func (c *Client) Upload(source io.Reader, destination io.Writer, size int) error {
	if err := c.connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	chunk := make([]byte, size)

	for {
		num, err := source.Read(chunk)
		if err == io.EOF {
			tot, err := destination.Write(chunk[:num])
			if err != nil {
				return err
			}

			if tot != len(chunk[:num]) {
				return fmt.Errorf("failed to write stream")
			}

			return nil
		}

		if err != nil {
			return err
		}

		tot, err := destination.Write(chunk[:num])
		if err != nil {
			return err
		}

		if tot != len(chunk[:num]) {
			return fmt.Errorf("failed to write stream")
		}
	}
}

// Download returns remote/destination file for reading.
func (c *Client) Download(filePath string) (io.ReadCloser, error) {
	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return c.sftpClient.OpenFile(filePath, (os.O_RDONLY))
}
func (sc *Client) ListFiles(remoteDir string) (err error) {
	fmt.Fprintf(os.Stdout, "Listing [%s] ...\n\n", remoteDir)

	files, err := sc.sftpClient.ReadDir(remoteDir)
	if err != nil {
		log.Printf("Unable to list remote dir: %v\n", err)
		return
	}

	for _, f := range files {
		var name, modTime, size string
		name = f.Name()
		modTime = f.ModTime().Format("2006-01-02 15:04:05")
		size = fmt.Sprintf("%12d", f.Size())

		if f.IsDir() {
			name = name + "/"
			modTime = ""
			size = "PRE"
		}
		// Output each file name and size in bytes
		fmt.Fprintf(os.Stdout, "%19s %12s %s\n", modTime, size, name)
	}

	return
}
func (c *Client) Remove(filePath string) error {
	if err := c.connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	return c.sftpClient.Remove(filePath)
}

// Info gets the details of a file. If the file was not found, an error is returned.
func (c *Client) Info(filePath string) (os.FileInfo, error) {
	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	info, err := c.sftpClient.Lstat(filePath)
	if err != nil {
		return nil, fmt.Errorf("file stats: %w", err)
	}

	return info, nil
}

// Close closes open connections.
func (c *Client) Close() {
	if c.sftpClient != nil {
		c.sftpClient.Close()
	}
	if c.sshClient != nil {
		c.sshClient.Close()
	}
}

// connect initialises a new SSH and SFTP client only if they were not
// initialised before at all and, they were initialised but the SSH
// connection was lost for any reason.
func (c *Client) connect() error {
	var (
		signer ssh.Signer
		err    error
	)
	if c.sshClient != nil {
		_, _, err := c.sshClient.SendRequest("keepalive", false, nil)
		if err == nil {
			return nil
		}
	}

	auth := ssh.Password(c.config.Password)
	if c.config.PrivateKey != "" {
		if c.config.KeyPsw != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(c.config.PrivateKey), []byte(c.config.KeyPsw))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(c.config.PrivateKey))
		}

		if err != nil {
			return fmt.Errorf("ssh parse private key: %w", err)
		}
		auth = ssh.PublicKeys(signer)
	}
	log.Println(signer.PublicKey())
	log.Println(auth)
	cfg := &ssh.ClientConfig{
		User: c.config.Username,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: func(string, net.Addr, ssh.PublicKey) error { return nil },
		Timeout:         c.config.Timeout,
		Config: ssh.Config{
			KeyExchanges: c.config.KeyExchanges,
		},
	}

	sshClient, err := ssh.Dial("tcp", c.config.Server, cfg)
	if err != nil {
		return fmt.Errorf("ssh dial: %w", err)
	}
	c.sshClient = sshClient

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("sftp new client: %w", err)
	}
	c.sftpClient = sftpClient

	return nil
}
func SftpUpload(host string, user string, pass string) {

	// Default SFTP port
	port := 22

	hostKey := getHostKey(host)

	fmt.Fprintf(os.Stdout, "Connecting to %s ...\n", host)

	var auths []ssh.AuthMethod

	// Try to use $SSH_AUTH_SOCK which contains the path of the unix file socket that the sshd agent uses
	// for communication with other processes.
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}

	// Use password authentication if provided
	if pass != "" {
		auths = append(auths, ssh.Password(pass))
	}

	// Initialize client configuration
	config := ssh.ClientConfig{
		User: user,
		Auth: auths,
		// Uncomment to ignore host key check
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	// Connect to server
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connecto to [%s]: %v\n", addr, err)
		os.Exit(1)
	}

	defer conn.Close()

	// Create new SFTP client
	sc, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start SFTP subsystem: %v\n", err)
		os.Exit(1)
	}
	defer sc.Close()
}
func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read known_hosts file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing %q: %v\n", fields[2], err)
				os.Exit(1)
			}
			break
		}
	}

	if hostKey == nil {
		fmt.Fprintf(os.Stderr, "No hostkey found for %s", host)
		os.Exit(1)
	}

	return hostKey
}
