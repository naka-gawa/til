package httpapi.authz

# 社員情報をインポート
import data.subordinates as subord

default allow = false

# 社員は自身の給与を参照できる
allow {
  some username
  input.method == "GET"
  input.path = ["finance", "salary", username]
  input.user == username
}

# マネージャーは部下の給与も参照できる
allow {
  some username
  input.method == "GET"
  input.path = ["finance", "salary", username]
  subord[input.user][_] == username
}
