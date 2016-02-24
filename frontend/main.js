import moment from 'moment'
import YouTube from 'react-youtube'

view SearchBox {
  <label>Search terms:</label>
  <br/>
  <input placeholder='cupcakes' onChange={view.props.onChange}/>
}

view DateBox {
  <label>{view.props.label} date:</label>
  <br/>
  <input placeholder='MM/DD/YYYY' onChange={view.props.onChange}/>
  $ = {
    marginTop: 10
  }
}

view SearchButton {
  <button onClick={view.props.onClick}>Search</button>
  $ = {
    marginTop: 10
  }
}
view CheckBox {
  <confirmationText> <input type='checkbox' onChange={view.props.onChange}/> localhost:5000 search </confirmationText>

}

view Video {
  <h5>{`${view.props.video.snippet.title}  — (${moment(view.props.video.snippet.publishedAt).format('MMMM DD, YYYY')})`}</h5>
  <YouTube
    opts={{height: '289', width: '360'}}
    videoId={view.props.video.id.videoId}
  />
  $ = {
    width: 360,
    outline: '1px solid black'
  }
}

view Main {
  <h1>tapestry</h1>

  results = []

  start = ""
  end = ""
  searchQuery = ""
  localhost = false

  async function search() {
    const startDate = moment(start, "MM-DD-YYYY")
    const endDate = moment(end, "MM-DD-YYYY")
    url = !localhost ? "http://tapestry-server.herokuapp.com/search?" : "//localhost:5000/search?"

    searchQuery.replace(" ", "%20")
    const test = `${url}q=${searchQuery}&after=${startDate.toISOString()}&before=${endDate.toISOString()}`
    results = await fetch.json(test)
    view.update()
  }

  <SearchBox onChange={e => searchQuery = e.target.value}/>
  <DateBox label={'Start'} onChange={e => start = e.target.value}/>
  <DateBox label={'End'} onChange={e => end = e.target.value}/>
  <SearchButton onClick={search}/>
  <CheckBox onChange={e => localhost = !localhost}/>
  // <CheckBox/>
  <videos>
    <Video repeat={results} video={_}/>
  </videos>

  $h1 = {
    color: 'dimgray',
  }

  $ = {
    padding: 45,
    paddingTop: 17
  }
}
