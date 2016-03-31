import _ from 'lodash'
import moment from 'moment'
import 'whatwg-fetch'
import 'moment-duration-format'
import React from 'react'
import ReactDOM from 'react-dom'
import classNames from 'classnames'

import FontAwesome from 'react-fontawesome'
import YouTube from 'react-youtube'

import './style.less'

 const setItem = (obj, x, xs) =>
   xs.map(_x => _x.id == x.id ? Object.assign(x, obj) : _x)

 const setAll = (obj, xs) => xs.map(x => Object.assign(x, obj))

 const browse = (c, cs) =>
   setItem({ isBrowsing: true }, c, setAll({ isBrowsing : false }, cs))

 const toggleSelected = (c, cs) =>
   setItem({ isSelected: !c.isSelected }, c, cs)

class DownloadButton extends React.Component {
  render() {
    const numClusters = this.props.numClusters
    const numVideos = this.props.numVideos

    return (
      <button disabled={numClusters == 0} onClick={this.props.onClick}>Download {numClusters} Clusters ({numVideos} videos)</button>
    )
  }

  // $button = {
  //   fontSize: 20,
  //   marginTop: 10,
  // }
}

class IconButton extends React.Component {
  render() {
    return <FontAwesome
      name={this.props.name}
      onClick={() => alert(this.props.name)}
      className='icon-button'
      />
  }
}

class VideoThumbnail extends React.Component {
  render() {
    if (!this.props.src) return <strong>invalid thumbnail</strong>
    return (
      <div className='video-thumbnail' onClick={this.props.onClick}>
        <img src={this.props.src}/>
        <div className='buttons'>
          <IconButton name='flag'/>
          <IconButton name='star'/>
          <IconButton name='trash'/>
        </div>
      </div>
    )
  }
}

class ClusterBrowser extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      index: 0,
      selectedIndex: null,
    }
  }

  render() {
    const step = 4
    const videos = this.props.videos
    const index = this.state.index
    const selectedIndex = this.state.selectedIndex
    const selectedVideo = videos[this.state.selectedIndex || index]

    const thumbnails = _.range(0, videos.length).map(i => {
      // swap video if its selected
      i = (i === selectedIndex) ? index : i
      return <VideoThumbnail
        key={i}
        onClick={e => this.setState({selectedIndex: i})}
        src={videos[i].thumbnail_url}
      />
    })

    return (
      <div className='cluster-browser'>
        <h3>Title/Date + Geographic Location {videos.length} videos</h3>
        <h5>{selectedVideo.title}</h5>
        <YouTube videoId={selectedVideo.youtube_id}/>
        {thumbnails.slice(index + 1, index + step)}
        <input
          className="scrubber"
          type="range"
          value={index}
          onChange={e => this.setState({
            index: Number(e.target.value),
            selectedIndex: null,
          })}
          min={0}
          step={step}
          max={videos.length - 1}
        />
        <div className="scrubber-label">
          {index + 1}-{index + step}/{videos.length}
        </div>
      </div>
    )
  }
}

class ClusterList extends React.Component {

  render() {
    return (
      <div>
        <strong>Next Clusters</strong>
      </div>
    )
  }
}

class Main extends React.Component {

  constructor(props) {
    super(props);
    this.state =  {
      clusters: [],
      index: 0,
    }
  }

  componentDidMount() {
    this.load()
  }

  load() {
    console.log('load')
    fetch('/api/videos')
    .then(response => response.json())
    .then(json => {
      const videos = json
      const clusters = _.map(_.groupBy(videos, 'tag'), (videos, tag) => {
        return {
          _id: tag || 'none',
          videos,
        }
      })
      this.setState({
        videos,
        clusters
      })
    })
    .catch(err => {
      console.log(err)
    })
  }

  render() {
    console.log(this.state);
    if (!this.state.clusters.length > 0) {
      return <strong>Loading...</strong>
    }
    return (
      <div className='main'>
        <ClusterList clusters={this.state.clusters}/>
        <ClusterBrowser videos={this.state.clusters[this.state.index].videos}/>
      </div>
    )
  }
}

ReactDOM.render(<Main/>, document.getElementById('root'))
