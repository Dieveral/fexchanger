<html>
<head>
	<title>Upload file</title>
</head>
<body>
<form enctype="multipart/form-data" action="http://localhost:8080/upload" method="post">
	<input type="file" name="uploadFiles" multiple />
	<input type="submit" name="upload" />
</form>
</body>
</html>