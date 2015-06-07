{{ define "content" }}

    <div class="container">
      <h2>TEST</h2>
      <form role="form" method="GET" action="/recipe">
        <div class="form-group">
          <input name="testValue" type="text" class="form-control" placeholder="Enter Sample Value">
        </div>
        <button type="submit" class="btn btn-default">OK</button>
      </form>
      <div class="form-group">
      </div>
    </div>
    
{{ end }}