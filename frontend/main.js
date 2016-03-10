import moment from 'moment'
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
  <img class='third' src='//placehold.it/250x250/0000cc'/>
  <img class='second' src='//placehold.it/250x250/00cc00'/>
  <img class='first' src='//placehold.it/250x250/cc0000'/>

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
    height: 150,
  }
}

view ClusterThumbnail {

  isSelected = false

  <input
    type='checkbox'
    defaultChecked={isSelected}
    onChange={e => isSelected = e.target.checked}
  />
  <StackedImages/>
  <detail class='title'>Title of Video</detail>
  <detail class='location'>Geographic Location</detail>
  <detail class='stats'>
    <detail class='count'>4 videos, </detail>
    <detail class='time'>1:06:40</detail>
  </detail>

  $ = {
    outline: '1px solid black',
    width: 250,
    fontSize: 12,
  }

  $StackedImages = {
    width: '92%',
    marginTop: -19
  }

  $input = {
    marginLeft: -20,
  }

  $title = {
    fontWeight: 'bold',
    fontSize: 14,
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
  <button>Download Selected Clusters</button>

  $button = {
    fontSize: 20,
    marginTop: 10,
  }
}

view Main {
  <h1>Most Recent Clusters</h1>

  <ClusterThumbnail/>

  <DownloadButton/>

  $ = {
    padding: 45,
    paddingTop: 17
  }
}
