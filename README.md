# Cookie in Go
Cookies are a way to store information at the client end. The client can be a browser, a mobile application, or anything which makes an HTTP request. Cookies are basically some files that are stored in the cache memory of your browser. When you are browsing any website which supports cookies will drop some kind of information related to your activities in the cookie. This information could be anything. Cookies in short store historical information about the user activities. This information is stored on the client’s computer. Since a cookie is stored in a file,  hence this information is not lost even when the user closes a browser window or restarts the computer. A cookie can also store the login information. In fact, login information such as tokens is generally stored in cookies only. Cookies are stored per domain. Cookies stored locally belonging to a particular domain are sent in each request to that domain. They are sent in each request as part of headers. So essentially cookie is nothing but a header.

# Fields in a Cookie

`Name` is the cookie name. It can contain any US-ASCII characters except ( ) < > @ , ; : \ " / [ ? ] = { } and space, tab and control characters. It is a mandatory field.

`Value` contains the data that you want to persist. It can contain any US-ASCII characters except , ; \ " and space, tab and control characters. It is a mandatory field.

`Path` attribute plays a major role in setting the scope of the cookies in conjunction with the domain. In addition to the domain, the URL path that the cookie is valid for can be specified. If the domain and path match, then the cookie will be sent in the request. Just as with the domain attribute, if the path attribute is set too loosely, then it could leave the application vulnerable to attacks by other applications on the same server. For example, if the path attribute was set to the web server root /, then the application cookies will be sent to every application within the same domain (if multiple application reside under the same server). A couple of examples for multiple applications under the same server:

    path=/bank
    path=/private
    path=/docs
    path=/docs/admin

`Domain` attribute is used to compare the cookie’s domain against the domain of the server for which the HTTP request is being made. If the domain matches or if it is a subdomain, then the path attribute will be checked next.

Note that only hosts that belong to the specified domain can set a cookie for that domain. Additionally, the domain attribute cannot be a top level domain (such as .gov or .com) to prevent servers from setting arbitrary cookies for another domain (such as setting a cookie for owasp.org). If the domain attribute is not set, then the hostname of the server that generated the cookie is used as the default value of the domain.

For example, if a cookie is set by an application at app.mydomain.com with no domain attribute set, then the cookie would be resubmitted for all subsequent requests for app.mydomain.com, but not its subdomains (such as hacker.app.mydomain.com), or to otherapp.mydomain.com. (However, older versions of Edge/IE behave differently, and do send these cookies to subdomains.) If a developer wanted to loosen this restriction, then they could set the domain attribute to mydomain.com. In this case the cookie would be sent to all requests for app.mydomain.com and mydomain.com subdomains, such as hacker.app.mydomain.com, and even bank.mydomain.com. If there was a vulnerable server on a subdomain (for example, otherapp.mydomain.com) and the domain attribute has been set too loosely (for example, mydomain.com), then the vulnerable server could be used to harvest cookies (such as session tokens) across the full scope of mydomain.com.  

`Expires` set persistent cookies
limit lifespan if a session lives for too long
remove a cookie forcefully by setting it to a past date
When a cookie passes its expiry date, it will no longer be sent with browser requests, and instead will be deleted. The date value is a HTTP timestamp.

`Secure` means Cookie will only be sent when the request is made with HTTPS.
The Secure attribute tells the browser to only send the cookie if the request is being sent over a secure channel such as HTTPS. This will help protect the cookie from being passed in unencrypted requests. If the application can be accessed over both HTTP and HTTPS, an attacker could be able to redirect the user to send their cookie as part of non-protected requests.

`HttpOnly` attribute is used to help prevent attacks such as session leakage, since it does not allow the cookie to be accessed via a client-side script such as JavaScript.

This doesn’t limit the whole attack surface of XSS attacks, as an attacker could still send request in place of the user, but limits immensely the reach of XSS attack vectors.


`SameSite` allows a server to define a cookie attribute making it impossible for the browser to send this cookie along with cross-site requests. The main goal is to mitigate the risk of cross-origin information leakage, and provide some protection against cross-site request forgery attacks.

SameSite attribute can be used to assert whether a cookie should be sent along with cross-site requests. This feature allows the server to mitigate the risk of cross-origin information leakage. In some cases, it is used too as a risk reduction (or defense in depth mechanism) strategy to prevent cross-site request forgery attacks. This attribute can be configured in three different modes:

    Strict
    Lax
    None

`Strict` value is the most restrictive usage of SameSite, allowing the browser to send the cookie only to first-party context without top-level navigation. In other words, the data associated with the cookie will only be sent on requests matching the current site shown on the browser URL bar. The cookie will not be sent on requests generated by third-party sites. This value is especially recommended for actions performed at the same domain. However, it can have some limitations with some session management systems negatively affecting the user navigation experience. Since the browser would not send the cookie on any requests generated from a third-party domain or email, the user would be required to sign in again even if they already have an authenticated session.
Lax Value

`Lax` value is less restrictive than Strict. The cookie will be sent if the URL equals the cookie’s domain (first-party) even if the link is coming from a third-party domain. This value is considered by most browsers the default behavior since it provides a better user experience than the Strict value. It doesn’t trigger for assets, such as images, where cookies might not be needed to access them.
None Value

`None value` specifies that the browser will send the cookie in all contexts, including cross-site requests (the normal behavior before the implementation of SameSite). If Samesite=None is set, then the Secure attribute must be set, otherwise modern browsers will ignore the SameSite attribute, e.g. SameSite=None; Secure.

```go 
const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)
```

# Cookie struct in Go

A Cookie represents an HTTP cookie as sent in the Set-Cookie header of an HTTP response or the Cookie header of an HTTP request. 

```go 
type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

## Intialization of a new cookie instance

```go 
    cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "Hello Shubham Mishra!",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
}
```

## Setting a cookie

```http.SetCookie(w, &cookie)```

## To get a cookie

```cookie, err := r.Cookie("exampleCookie")```

## To clear a cookie

Must note 2 things at least while clearing a cookie first its `Name` second `MaxAge<0`

`Name` should be same name of cookie which you want to clear.

```go 
c := &http.Cookie{
		Name:     "exampleCookie",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
    }

http.SetCookie(w, c)
```

# Operation on a cookie 

consider a cookie instance `c` where c is a cookie.

`c.String() string`
String() method returns the serialization of the cookie for use in a Cookie header (if only Name and Value are set) or a Set-Cookie response header (if other fields are set). If c is nil or c.Name is invalid, the empty string is returned. 

`c.Valid() error`
Valid reports whether the cookie is valid. 

Of course all struct realated terms like `Name`,`Value`,`Path`,`Domain` etc... you can access using cookie instance `c`.

# Refernce

```https://tools.ietf.org/html/rfc6265 ```

```https://pkg.go.dev/net/http#Cookie.HttpOnly```