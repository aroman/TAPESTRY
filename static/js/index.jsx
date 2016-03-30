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
    return (
      <div className='video-thumbnail'>
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
    this.state =  {
      index: 0,
    }
  }

  render() {
    const videos = this.props.videos
    const root = videos[0]
    const index = this.state.index
    return (
      <div className='cluster-browser'>
        <h3>{root.title} Date + Geographic Location {videos.length} videos</h3>
        <h5>{videos[index].title}</h5>
        <YouTube videoId={videos[index].youtube_id}/>
        <VideoThumbnail src={videos[index + 1].thumbnail_url} />
        <VideoThumbnail src={videos[index + 2].thumbnail_url}/>
        <VideoThumbnail src={videos[index + 3].thumbnail_url}/>
        <input
          className="scrubber"
          type="range"
          value={index}
          onChange={e => this.setState({index: Number(e.target.value)})}
          min={0}
          max={this.props.videos.length - 4}
        />
        <div className="scrubber-label">
          {this.state.index + 1}-{this.state.index + 4}/{videos.length}
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
