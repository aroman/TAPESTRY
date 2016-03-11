import moment from 'moment'
import _ from 'underscore'
import YouTube from 'react-youtube'

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

view StackedImages {
  prop images

  <img class='third' if={images[2]} src={images[2].thumbnail_url}/>
  <img class='second' if={images[1]} src={images[1].thumbnail_url}/>
  <img class='first' if={images[0]} src={images[0].thumbnail_url}/>

  offset = 10

  $ = {
    position: 'relative',
    marginBottom: offset * 3,
  }

  $third = { top: 0, left: 0 }

  $second = { top: offset, left: offset }

  $first = {
    top: offset * 2,
    left: offset * 2,
    position: "relative",
  }

  $img = {
    position: 'absolute',
    width: '100%',
    height: 90,
  }
}

view ClusterThumbnail {

  isSelected = false

  on.click(view.props.onClick)

  <input
    type='checkbox'
    defaultChecked={isSelected}
    onChange={e => isSelected = e.target.checked}
  />
  <StackedImages images={view.props.videos.slice(0,3)}/>
  <detail class='title'>{view.props.videos[0].title}</detail>
  <br/>
  <detail class='location'>{
    `${view.props.videos[0].latitutde},${view.props.videos[0].longitude}`
  }</detail>
  <detail class='stats'>
    <detail class='count'>{view.props.videos.length} videos, </detail>
    <detail class='time'>1:06:40</detail>
    <detail class='tag'>tag: {view.props.id}</detail>
  </detail>

  $ = {
    // outline: '1px solid black',
    width: 250,
    fontSize: 12,
    margin: 25,
    flex: 1,
  }

  $StackedImages = {
    width: 120,
    marginBottom: 30,
  }

  $input = {
    marginLeft: -20,
  }

  $title = {
    fontWeight: 'bold',
    fontSize: 12,
  }

  $location = {
    fontSize: 12,
  }

  $count = {
    display: 'inline'
  }

  $time = {
    fontWeight: 'bold',
    display: 'inline',
  }

  $stats = {
    float: 'right',
    marginTop: -32,
    marginRight: 5,
  }
}

view DownloadButton {
  <button onClick={view.props.onClick}>Download Selected Clusters</button>

  $button = {
    fontSize: 20,
    marginTop: 10,
  }
}

view Main {
  videos = []
  clusters = []
  clusterIndex = null

  function load() {
    console.log('load')
    fetch('http://localhost:8000/')
    .then(response => response.json())
    .then(json => {
      videos = json
      clusters = _.map(_.groupBy(videos, 'tag'), (videos, tag) => {
        return {
          _id: tag || 'none',
          videos,
        }
      })
      view.update()
    })
    .catch(err => {
      console.log(err)
    })
  }

  load()

  <h1>Most Recent Clusters</h1>

  <viewer>
    <YouTube
      if={clusterIndex !== null}
      opts={{height: '289', width: '360'}}
      videoId={clusters[clusterIndex].videos[0].youtube_id}
    />
  </viewer>

  <grid>
    <ClusterThumbnail
      repeat={clusters}
      id={_._id}
      videos={_.videos}
      onClick={() => {clusterIndex = clusters.indexOf(_); view.update()}}
    />
  </grid>

  <DownloadButton onClick={() => alert('downloading...')}/>

  $ = {
    padding: 45,
    paddingTop: 20
  }

  $grid = {
    display: 'flex',
    flexFlow: 'row wrap',
  }
}
