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

const LABEL_TYPE = {
  Unmarked: 'unmarked',
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
      selectedIndex: null,
    }
  }

  render() {
    const step = 4
    const cluster = this.props.cluster
    const videos = cluster.videos
    const index = this.state.index
    const selectedIndex = this.state.selectedIndex
    const selectedVideo = videos[this.state.selectedIndex || index]

    const thumbnails = _.range(0, videos.length).map(i => {
      // swap video if its selected
      i = (i === selectedIndex) ? index : i
      return <img
        key={i}
        className='video-thumbnail'
        src={videos[i].thumbnail_url}
        onClick={e => this.setState({selectedIndex: i})}
      />
    })

    const filterButtons = _.map(LABEL_TYPE, (value, key) => {
      return <IconButton
        name={value}
        isActive={cluster.label == LABEL_TYPE[key]}
        onClick={() => this.props.onSetLabel(cluster, LABEL_TYPE[key])}
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
            selectedIndex: null,
          })}
          min={0}x
          step={step}
          max={videos.length - 1}
        />
        <div className='controls'>
          <div className='scrubber-label'>
            {index + 1}-{Math.min(videos.length, index + step)}/{videos.length}
          </div>
          {filterButtons}
          <button onClick={this.props.onNext}>Next</button>
        </div>
      </div>
    )
  }
}

class ClusterList extends React.Component {

  constructor(props) {
    super(props)
    this.state = {
      filterMode: LABEL_TYPE.Unmarked,
    }
  }

  filterClicked(mode) {
    if (this.state.filterMode == mode) {
      this.setState({filterMode: LABEL_TYPE.Unmarked})
      return;
    }
    this.setState({
      filterMode: mode,
    })
  }

  render() {
    const clusters = this.props.clusters
    const selectedIndex = this.props.selectedIndex

    let filcos = clusters.filter(cluster => cluster.label == this.state.filterMode)

    const thumbnails = _.range(filcos.length).map(i => {
      const root = filcos[i].videos[0]

      return (
        <div key={filcos[i].id} onClick={() => this.props.onChange(i)} className={classNames('sidebar-item', {
            selected: i == selectedIndex
          })}>
          <img
            className='sidebar-item-img'
            src={root.thumbnail_url}
          />
          <span className='sidebar-item-title'>{root.title}</span>
        </div>
      )
    })

    const filterButtons = _.map(LABEL_TYPE, (value, key) => {
      return <IconButton
        name={value}
        isDark={true}
        isActive={this.state.filterMode == LABEL_TYPE[key]}
        onClick={() => this.filterClicked(LABEL_TYPE[key])}
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
        <div className='videos'>
          {thumbnails}
        </div>
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

  onSetLabel(cluster, label) {
    if (cluster.label == label) {
      label = LABEL_TYPE.Unmarked;
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
    console.log('load')
    fetch('/api/videos')
    .then(response => response.json())
    .then(json => {
      const videos = json
      const clusters = _.map(_.groupBy(videos, 'tag'), (videos, tag) => {
        return {
          id: tag || 'none',
          label: localStorage[tag] || LABEL_TYPE.Unmarked,
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
    if (!this.state.clusters.length > 0) {
      return (
        <div className='loading'>
          <img src='/static/img/loading.gif'/>
        </div>
      )
    }
    return (
      <div className='main'>
        <ClusterBrowser
          cluster={this.state.clusters[this.state.index]}
          onSetLabel={this.onSetLabel.bind(this)}
          onNext={() => this.setState({index: this.state.index + 1})}
        />
        <ClusterList
          clusters={this.state.clusters}
          selectedIndex={this.state.index}
          onChange={clusterIndex => this.setState({index: clusterIndex})}
        />
      </div>
    )
  }
}

ReactDOM.render(<Main/>, document.getElementById('root'))
