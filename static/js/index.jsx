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

const LABEL = {
  Unmarked: '',
  Star: 'star',
  Flag: 'flag',
  Trash: 'trash',
}

class IconButton extends React.Component {

  render() {
    let color = 'white'

    if (this.props.isDark) {
      color = 'black'
    }

    if (this.props.isActive) {
      color = 'red'
    }

    return <FontAwesome
      name={this.props.name}
      style={{color: color}}
      size='2x'
      onClick={this.props.onClick}
      className='icon-button'
      />
  }
}

class ClusterBrowser extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      index: 0,
      selectedVideoId: null,
    }
  }

  render() {
    const step = 4
    const videos = this.props.cluster.videos
    const index = this.state.index
    const selectedVideo = _.find(videos, v => v.id === this.state.selectedVideoId) || this.props.cluster.videos[0]

    const thumbnails = this.props.cluster.videos.map((video, i) => {
      // swap video if its selected
      if (video.id == selectedVideo.id) {
        video = videos[index]
      }
      return <img
        key={video.id}
        className='video-thumbnail'
        src={video.thumbnail_url}
        onClick={e => this.setState({selectedVideoId: video.id})}
      />
    })

    const filterButtons = _.map(LABEL, (value, key) => {
      return <IconButton
        key={key}
        name={value}
        isActive={this.props.cluster.label == LABEL[key]}
        onClick={() => this.props.onSetLabel(cluster, LABEL[key])}
      />
    })

    return (
      <div className='cluster-browser'>
        <h5>{selectedVideo.title}</h5>
        <YouTube opts={{width:560, height:315}} videoId={selectedVideo.youtube_id}/>
        <div className='thumbnails'>
          {thumbnails.slice(index + 1, index + step)}
        </div>
        <input
          className="scrubber"
          type="range"
          value={index}
          onChange={e => this.setState({
            index: Number(e.target.value),
          })}
          min={0}
          step={step}
          max={videos.length - 1}
        />
        <div className='controls'>
          <div className='scrubber-label'>
            {index + 1}-{Math.min(videos.length, index + step)}/{videos.length}
          </div>
          {filterButtons}
        </div>
      </div>
    )
  }
}

class SidebarItemList extends React.Component {

  render() {
    const items = this.props.clusters.map(cluster => {
      const root = cluster.videos[0]

      return (
        <div key={cluster.id} onClick={() => this.props.onChange(cluster)} className={classNames('sidebar-item', {
            selected: cluster.id == this.props.selectedCluster.id
          })}>
          <img
            className='sidebar-item-img'
            src={root.thumbnail_url}
          />
          <span className='sidebar-item-title'>{root.title}</span>
        </div>
      )
    })
    return (
      <div className='videos'>
        {items}
      </div>
    )
  }
}

class Sidebar extends React.Component {

  constructor(props) {
    super(props)
    this.state = {
      filterMode: LABEL.Unmarked,
    }
  }

  filterClicked(mode) {
    if (this.state.filterMode == mode) {
      this.setState({filterMode: LABEL.Unmarked})
      return;
    }
    this.setState({
      filterMode: mode,
    })
  }

  render() {
    const results = this.props.clusters.filter(cluster => cluster.label === this.state.filterMode)

    const filterButtons = _.map(LABEL, (value, key) => {
      return <IconButton
        key={key}
        name={value}
        isDark={true}
        isActive={this.state.filterMode == LABEL[key]}
        onClick={() => this.filterClicked(LABEL[key])}
      />
    })

    return (
      <div className='cluster-list'>
        <div className='list-header'>
          <h3>Next Clusters</h3>
          <div className='filters'>
            <span>Filters</span>
            {filterButtons}
          </div>
        </div>
        <SidebarItemList
          clusters={results}
          selectedCluster={this.props.selectedCluster}
          onChange={this.props.onChange}
        />
      </div>
    )
  }
}

class Main extends React.Component {

  constructor(props) {
    super(props);
    this.state =  {
      clusters: [],
      selectedClusterId: null,
    }
  }

  onSetLabel(cluster, label) {
    if (cluster.label == label) {
      label = LABEL.Unmarked;
    }
    this.setState({
      clusters: setItem({label: label}, cluster, this.state.clusters),
    })
    localStorage.setItem(cluster.id, label)
  }

  componentDidMount() {
    this.load()
  }

  load() {
    fetch('/api/clusters')
    .then(response => response.json())
    .then(clusters => {
      this.setState({
        selectedClusterId: clusters[0].id,
        clusters,
      })
    })
    .catch(err => {
      console.log(err)
    })
  }

  render() {
    if (!this.state.clusters.length > 0) {
      return (
        <div className='loading'>
          <img src='/static/img/loading.gif'/>
        </div>
      )
    }
    const selectedCluster = _.find(this.state.clusters, c => c.id == this.state.selectedClusterId) || this.state.clusters[0]
    return (
      <div className='main'>
        <ClusterBrowser
          cluster={selectedCluster}
          onSetLabel={this.onSetLabel.bind(this)}
        />
        <Sidebar
          clusters={this.state.clusters}
          selectedCluster={selectedCluster}
          onChange={cluster => this.setState({ selectedClusterId: cluster.id })}
        />
      </div>
    )
  }
}

ReactDOM.render(<Main/>, document.getElementById('root'))
