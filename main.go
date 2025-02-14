package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"golang.org/x/net/http2"
	"net/http/httptrace"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
		"mime/multipart"
)

func main() {

	// cert := "./server.crt"
	// key := "./server.key"

	f, err := os.Create("/tmp/trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	trace.Start(f)
	defer trace.Stop()

	cert := []byte(`-----BEGIN CERTIFICATE-----
MIIDETCCAfkCFBhM16fLpiiNMIN/z95CnDi1u6qiMA0GCSqGSIb3DQEBCwUAMEUx
CzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRl
cm5ldCBXaWRnaXRzIFB0eSBMdGQwHhcNMjMwOTI2MDAwNzI2WhcNMjMxMDI2MDAw
NzI2WjBFMQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UE
CgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAs2ykz5Zv4EDxWMqGX6YRBeSSsoc+/12L4VkNhc/lx5UlkWIP
sAhJ8LMZ+KHXxZ8F11tBfy+vt6SlqpD/Gr+Cu0wpwINIawPbNBUW4T6h0Ch5EuVa
0R+9tDRYcCw61Kvv9CN4BhN4ngy2j1DvYZoL6Mo8LYXJmhOns2vdT0/MQHs5Xgy3
LB6ozof1XCpv31JB/dQmZdrGmk6wFwcv5jDO30aWZ03otWFfG9m9hSHY/mwQigmr
Un+DE56kkdw8XmK90wRyVjgVo0VZIbXi6k8DBVOvxCkSzsNUVz0r0B4dBdAzHafA
c2eWaeJnlMPFrrCG1zUrZYz6C/Gs0tuy6TbwZQIDAQABMA0GCSqGSIb3DQEBCwUA
A4IBAQAagO82txNlb0iH1YWQGiCVbHgwKS7Zf3h8A6sO6PsOHAbid35Cef/rXBpR
MfVarZcl/ku4c+dibTR4BrMJ3yip42/ZtJz5faBzCt4Bk+tsXVK+BGuTOslcGqAm
eWpiwlO4QBCeBcdHN5UdrymUl5RU8OSTubjV1jebUkoXb7zskKq0M3HpEmv31mke
pSLB1F3TqJ7TwuiQAX3B5MsfSq8dywdra4Xxa1/2s/yxE9DumeRj1mQ9ZKuBEA6G
QwWzutxUwe0PZVMLOzAREqvbyPiOfoKv28i516KdxAeyE9YEeVRW3ZoSuEjNZBfu
l6t75zuacx5cSxoqL4C3PsP79ghP
-----END CERTIFICATE-----`)

	key := []byte(`-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCzbKTPlm/gQPFY
yoZfphEF5JKyhz7/XYvhWQ2Fz+XHlSWRYg+wCEnwsxn4odfFnwXXW0F/L6+3pKWq
kP8av4K7TCnAg0hrA9s0FRbhPqHQKHkS5VrRH720NFhwLDrUq+/0I3gGE3ieDLaP
UO9hmgvoyjwthcmaE6eza91PT8xAezleDLcsHqjOh/VcKm/fUkH91CZl2saaTrAX
By/mMM7fRpZnTei1YV8b2b2FIdj+bBCKCatSf4MTnqSR3DxeYr3TBHJWOBWjRVkh
teLqTwMFU6/EKRLOw1RXPSvQHh0F0DMdp8BzZ5Zp4meUw8WusIbXNStljPoL8azS
27LpNvBlAgMBAAECggEBAIaoFfTZoXL7KmaJ8FkeiA4VImesCwKaEV35eKWhhoDI
FZ7LqpxhTCTCNc6coBRnqXmbukca4dSAU/jzwxEvA41PttINdS8jgMLQxRAc8lIH
3f3pdcQW/0ecNXdSKoDr3TUr3Wsp988eGXDrvAxTEXcPOCSuEMR61dRZMonQiKkm
zpIF9+EHg1bXwCyOe92qu9j1z5ahKBbkSw4jbSmndVulZr0bjazm329CQ7NIINx4
k3m0jU0vPj7Ae63sSy8ppkCW1c+DBkJWHT+hWSM8IUoL48ux/reexVfmdiAEj65v
P4B6b3Bavy0cNihfJXuD3yFCy75NF1q9otqb/dtqtq0CgYEA6/1reM5O/iBzuBp7
T6f61Sdx8pI6Lg/hXaY/y3ThfBy2i8sUlW7Fr1gg9QrTGWCnET3l4t9cqqcsSxXM
mhS7ztT4vj3Xh8ZaJCLF9/NP4dOF3a2BwpzztvQA+uWMMuvPPPeJUxfBIbkQRfw4
VRhSZtA8Vro0jxLhiaoUV409QIcCgYEAwqNeY1IhJtmYo8+l2T9WlaEdJaF642V/
KMddM6BCbPXriQnQxDPs1V5qmKMioFety1tbzCWEywQXUI5jL8i5UJzy8SbA+EN7
v19+hsNw95i35LwAgaXTBi7+w3C/RV+nwDvAIFfVh+RgwLm//mpxM+NhtqMLZgDO
0ZfN0EtVHrMCgYEAzNtZvy4A3pPXQHskjlp44S5zuSn8/V1lusEF4h1dXDsksCej
f7EARJ2RRoaWCBKiFNKRzCcvryidx7Rtv1e+TOcN4z+V6NRkDbht+DcsIhJXw77J
xOIwoqIgc6xuzwlrBjav7ATT7+1R1h8D28LYTU35AyRzWbv/M10EeZgvEuMCgYEA
sSnM73DBVkB37JtEZbN+V2I4SplqZvNqxC57ilNUQyOxCaKSkRSGfUyckGTlTreg
LNryeqan8ba+lgeVBtiYvWSuyYwx14htais82uOovuYWdLAfceMDarx6DmFg6H+C
DLsWTRpF9ZSN8L7ioGw4sDdHKNvNs0IG/sZDg1Yem3sCgYBzib9ifTOmTjE37C4x
id/w4CQ/jP4IB5MT69d3zBbqxyche0haUnIiZUneIF1AL7anAExTRbbMGivYflYo
kZ6eZuoK8diihUuFRv8Cplfr5yfZYDtEm1LPaxJiT2470zpLxD3Mks+obqjlHPLc
6M4C1VuJryQriqUW9ITvkHk6jw==
-----END PRIVATE KEY-----`)

	tr := &http.Transport{
		MaxIdleConns:        100000,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     time.Minute * time.Duration(1),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	//	ForceAttemptHTTP2: true,
	}
	 err = http2.ConfigureTransport(tr)
	 if err != nil {
			 panic(err)
	 }

	cnt := runtime.GOMAXPROCS(0)
	fmt.Println("cpu count ", cnt)

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {

			if !connInfo.Reused {
				fmt.Printf("new connection %v\n", connInfo.Conn.LocalAddr().String())
			}
		},
	}

	data := make([]byte, 9*1024*1024) // 9MB
	rand.Seed(time.Now().UnixNano())
	rand.Read(data)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(w, bytes.NewReader(data))
		if err != nil {
			panic(fmt.Sprintf("Internal Server Error %v", http.StatusInternalServerError))
		}
	})

	go func() {

		cert, err := tls.X509KeyPair(cert, key)
		if err != nil {
			panic(err)
		}
		server := &http.Server{
			Addr: ":12346",
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
			Handler: handler,
		}
		h2 := &http2.Server{}
		http2.ConfigureServer(server, h2)

		err = server.ListenAndServeTLS("", "")
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second * 1)

	startTime := time.Now()

			var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// 파일 추가
	part, err := writer.CreateFormFile("asdfb", "ab")
	if err != nil {
			panic(err)
	}

	if _, err := part.Write(data); err != nil {
			panic(err)
	}

	// 멀티파트 작성 종료
	if err := writer.Close(); err != nil {
			panic(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := "https://127.0.0.1:12346"

			req, err := http.NewRequest("POST", url, &body)
			if err != nil {
				panic(err)
			}
				req.Header.Set("Content-Type", writer.FormDataContentType())

			req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

			res, err := tr.RoundTrip(req)

			if err != nil {
				panic(err)
			}
			fmt.Println(res.Proto)

			_, err = io.ReadAll(res.Body)
			defer res.Body.Close()

			if err != nil {
				panic(err)
			}
		}()
	}

	wg.Wait()
	endTime := time.Now()
	fmt.Printf("time is %s", endTime.Sub(startTime))
}

