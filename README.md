# marionette_go
Fast dynamic site crawling based on chromedp (基于chromedp的动态网站抓取)

### deploy
- `docker-compose up -d`
### use marionette_go
#### ssr
- `get` http://127.0.0.1:6063/ssr?q=http://www.baidu.com
#### avaricious
-`post` http://127.0.0.1:6063/avaricious
- `json body`
```json 
{
    "url": "http://www.baidu.com",
    "timeout": 10
}
```
### Load Balancing
`docker-compose scale marionette=3`