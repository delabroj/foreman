package main

import "fmt"

func wrapHTML(body string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="theme-color" content="#000000">
		<style>
			pre {
			  white-space: pre-wrap;
				color: white;
				background-color: black;
				margin: grey 2px solid;
				padding: 10px;
			}
		</style>
    <title>Foreman</title>
  </head>
  <body>
    %s
  </body>
</html>
`, body)
}

var newPackageFrontend = wrapHTML(`
<h1>Upload new package</h1>
<script type="text/javascript">
  function checkForm(form) {
    form.myButton.disabled = true;
    form.myButton.value = "Please wait...";
    return true;
  }
</script>
<form method="post" enctype="multipart/form-data" onsubmit="return checkForm(this);">
  <div>
    <label for="package">Choose package to upload</label>
    <input type="file" id="package" name="package" accept=".zip"></br>
    <label for="hash">Package hash:</label>
    <textarea id="hash" name="hash" rows="3" cols="40"></textarea>
  </div>
  <div>
    <input type="submit" name="myButton" value="Submit"/>
  </div>
</form>
`)
