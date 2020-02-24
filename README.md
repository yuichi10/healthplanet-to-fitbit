# healthplanet-to-fitbit

## Reference
- healthplanet api
https://www.healthplanet.jp/apis/api.html

- healthplanet api oauth設定
https://www.healthplanet.jp/apis_account.do

- fitbit api
https://dev.fitbit.com/build/reference/web-api/  
https://dev.fitbit.com/build/reference/web-api/explore/  

- fitbit oauth設定
https://dev.fitbit.com/apps

## 設定ファイル
利用前に設定する項目  
userIDは設定しなくても大丈夫  
healthplanetと、fitbitのページに行ってoauthのアカウントを追加する。そのID/secretを追記する。  
設定ファイルは ~/.healthplanet_to_fitbit に設定する。  
```
healthplanet:
  userID:
  clientID:
  clientSecret:
fitbit:
  clientID:
  clientSecret:
```