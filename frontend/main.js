import moment from 'moment'
import YouTube from 'react-youtube'

view SearchBox {
  <label>Search terms:</label>
  <br/>
  <input placeholder='cupcakes' onChange={view.props.onChange}/>
}

view Location {
  <label>{view.props.label}</label>
  <br/>
  <input placeholder={view.props.placeholder} onChange={view.props.onChange}/>
  $ = {
    marginTop: 10
  }
  $input = {
    width: 100
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
  lat = ""
  long = ""
  radius = ""
  maxResults = ""
  localhost = false

  async function search() {
    const startDate = moment(start, "MM-DD-YYYY")
    const endDate = moment(end, "MM-DD-YYYY")
    url = !localhost ? "http://tapestry-server.herokuapp.com/search?" : "//localhost:5000/search?"

    searchQuery.replace(" ", "%20")
    const test = `${url}q=${searchQuery}
      &after=${startDate.toISOString()}
      &before=${endDate.toISOString()}
      &latitude=${lat}
      &longitude=${long}
      &radius=${radius}
      &maxResults=${maxResults}`
    results = await fetch.json(test)
    view.update()
  }

  <SearchBox onChange={e => searchQuery = e.target.value}/>
  <Location label={'Start'} placeholder={'MM/DD/YYYY'} onChange={e => start = e.target.value}/>
  <Location label={'End'} placeholder={'MM/DD/YYYY'} onChange={e => end = e.target.value}/>
  <Location label={'Latitude:'} onChange={e => lat = e.target.value}/>
  <Location label={'Longitude:'} onChange={e => long = e.target.value}/>
  <Location label={'Radius:'} placeholder={'1500m, 5km, 10000ft, and 0.75mi'} onChange={e => radius = e.target.value}/>
  <Location label={'Results Count:'} placeholder={'Between 0 and 50'} onChange={e => maxResults = e.target.value}/>
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
