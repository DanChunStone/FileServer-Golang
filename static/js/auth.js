function queryParams() {
  var username = localStorage.getItem("username");
  var token = localStorage.getItem("token");
    // console.log("username:"+username);
    // console.log("token:"+token);
  return 'username=' + username + '&token=' + encodeURI(token);
}