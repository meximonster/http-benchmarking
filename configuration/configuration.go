package configuration

import (
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var appConfig *Config

type AdminConfig struct {
}

type Config struct {
	Ip            string
	Port          string
	Threads       int
	Requests      int
	Endpoint      string
	ClientTimeout int
	Method        string
	Frequency     string
	Buckets       string
	Verbose       bool
	UUIDParam     bool
}

func Load() error {

	var c Config

	err := godotenv.Load()
	if err != nil {
		return err
	}

	c.Ip = os.Getenv("ip")
	if c.Ip == "" {
		c.Ip = "localhost"
	}

	c.Port = os.Getenv("port")
	if c.Port == "" {
		c.Port = "8080"
	}

	c.Endpoint = os.Getenv("endpoint")
	if c.Endpoint == "" {
		c.Endpoint = "https://www.google.gr"
	}

	timeout := os.Getenv("client_timeout")
	if timeout == "" {
		timeout = "10"
	}
	clTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		return err
	}
	c.ClientTimeout = clTimeout

	c.Method = os.Getenv("method")
	if c.Method == "" {
		c.Method = "GET"
	}

	c.Frequency = os.Getenv("frequency")
	if c.Frequency == "" {
		c.Frequency = "10s"
	}

	c.Buckets = os.Getenv("buckets")
	if c.Buckets == "" {
		c.Buckets = "50,100,500,1000"
	}

	thrds := os.Getenv("threads")
	if thrds == "" {
		thrds = "1"
	}
	t, err := strconv.Atoi(thrds)
	if err != nil {
		return err
	}
	c.Threads = t

	reqs := os.Getenv("requests")
	if reqs == "" {
		reqs = "3"
	}
	r, err := strconv.Atoi(reqs)
	if err != nil {
		return err
	}
	c.Requests = r

	vrb := os.Getenv("verbose")
	if vrb == "" {
		vrb = "false"
	}
	b, err := strconv.ParseBool(vrb)
	if err != nil {
		return err
	}
	c.Verbose = b

	param := os.Getenv("uuidParam")
	if param == "" {
		param = "false"
	}
	p, err := strconv.ParseBool(param)
	if err != nil {
		return err
	}
	c.UUIDParam = p

	appConfig = &c

	return nil
}

func Read() *Config {
	return appConfig
}

func (c *Config) Validate() ([]byte, time.Duration, []float64, error) {
	// Url validation
	_, err := url.ParseRequestURI(c.Endpoint)
	if err != nil {
		return nil, -1, nil, err
	}
	u, err := url.Parse(c.Endpoint)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, -1, nil, err
	}
	// Frequency validation
	freq, err := time.ParseDuration(c.Frequency)
	if err != nil {
		return nil, -1, nil, err
	}

	// Method validation
	var jbody []byte
	switch c.Method {
	case "GET":
		jbody = nil
	case "POST", "PUT", "PATCH":
		jsonFile, err := os.Open("req.json")
		if err != nil {
			return nil, -1, nil, err
		}
		defer jsonFile.Close()
		b, err := io.ReadAll(jsonFile)
		if err != nil {
			return nil, -1, nil, err
		}
		jbody = b
	default:
		log.Fatal("invalid http method")
	}

	// Bucket validation

	bucketArr := strings.Split(c.Buckets, ",")
	var bucketFloats []float64
	for _, b := range bucketArr {
		f, err := strconv.ParseFloat(b, 64)
		if err != nil {
			return nil, -1, nil, err
		}
		bucketFloats = append(bucketFloats, f)
	}
	return jbody, freq, bucketFloats, nil
}
