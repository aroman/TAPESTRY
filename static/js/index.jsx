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
    console.log(this.props)
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

    return (
      <div className='cluster-browser'>
        <h5>{selectedVideo.title}</h5>
        <YouTube opts={{width:560, height:315}} videoId={selectedVideo.youtube_id}/>
        {thumbnails.slice(index + 1, index + step)}
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
            {index + 1}-{index + step}/{videos.length}
          </div>
          <IconButton
            name='flag'
            isActive={cluster.label == LABEL_TYPE.Flag}
            onClick={() => this.props.onLabel(cluster, LABEL_TYPE.Flag)}
          />
          <IconButton
            name='star'
            isActive={cluster.label == LABEL_TYPE.Star}
            onClick={() => this.props.onLabel(cluster, LABEL_TYPE.Star)}
          />
          <IconButton
            name='trash'
            isActive={cluster.label == LABEL_TYPE.Trash}
            onClick={() => this.props.onLabel(cluster, LABEL_TYPE.Trash)}
          />
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

    return (
      <div className='cluster-list'>
        <div className='list-header'>
          <h3>Next Clusters</h3>
          <div className='filters'>
            <span>Filters</span>
            <IconButton
              name='flag'
              isDark={true}
              isActive={this.state.filterMode == LABEL_TYPE.Flag}
              onClick={() => this.filterClicked(LABEL_TYPE.Flag)}
            />
            <IconButton
              name='star'
              isDark={true}
              isActive={this.state.filterMode == LABEL_TYPE.Star}
              onClick={() => this.filterClicked(LABEL_TYPE.Star)}
            />
            <IconButton
              name='trash'
              isDark={true}
              isActive={this.state.filterMode == LABEL_TYPE.Trash}
              onClick={() => this.filterClicked(LABEL_TYPE.Trash)}
            />
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

  onClickLabel(cluster, label) {
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
          onLabel={this.onClickLabel.bind(this)}
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
