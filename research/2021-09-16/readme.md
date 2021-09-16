# September 16 2021

https://github.com/porschuetz/dev-budzillus.shalomsalon.de/blob/master/wordpress/wp-content/plugins/wp-bandcamp/components/bandcampClient.php

## /api/album/1/info

<http://bandcamp.com/api/album/1/info?key=thrjozkaskhjastaurrtygitylpt&album_id=79940049>

## /api/band/1/info

- <http://bandcamp.com/api/band/1/info?key=thrjozkaskhjastaurrtygitylpt&band_id=2853020814>
- <http://bandcamp.com/api/band/1/info?key=thrjozkaskhjastaurrtygitylpt&band_url=duststoredigital.com>

## /api/mobile/24/tralbum\_details

<http://bandcamp.com/api/mobile/24/tralbum_details?tralbum_type=t&band_id=2853020814&tralbum_id=714939036>

## /api/url/1/info

http://bandcamp.com/api/url/1/info?key=thrjozkaskhjastaurrtygitylpt&url=https://duststoredigital.com/track/trojan-horus-part-1

## /login\_cb

This doesnt work, as it requires Captcha

## /oauth\_login

This works:

~~~
POST /oauth_login HTTP/1.1
host: bandcamp.com
x-bandcamp-dm: 8f38339869c3003e9f1c8b1c13fe48530f74e3c6

client_id=134
client_secret=1myK12VeCL3dWl9o%2FncV2VyUUbOJuNPVJK6bZZJxHvk%3D
grant_type=password
password=PASSWORD
username=4095486538
username_is_user_id=1
~~~

We can get `x-bandcamp-dm` from Android, but its only good for three minutes. I
found an implementation online, but it seems BandCamp has changed the algorithm:

https://github.com/the-eater/camp-collective/issues/5

## /oauth\_token

We can try this:

~~~
POST /oauth_token HTTP/1.1
host: bandcamp.com

client_id=134&
client_secret=1myK12VeCL3dWl9o%2FncV2VyUUbOJuNPVJK6bZZJxHvk%3D&
grant_type=client_credentials
~~~

Result:

~~~
Only third-party clients can use client_credentials
~~~