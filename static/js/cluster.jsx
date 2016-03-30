import React from 'react'

class StackedImages extends React.Component {
  render() {
    return (
      <div>
        <img class='third' src={images[2].thumbnail_url}/>
        <img class='second' src={images[1].thumbnail_url}/>
        <img class='first' src={images[0].thumbnail_url}/>
      </div>
    )
  }
}

class ClusterExplorer extends React.Component {

  render() {
    let cluster = this.props.cluster


    let index = 0

    const prev = () => index = max(index-1, 0)
    const next = () => index = min(index+1, cluster.videos.length)


    return <strong>{'ClusterExplorer'}</strong>
    // return (
      // <div>
      //   <button onClick={prev}>{'<'}</button>
      //   <button onClick={next}>{'>'}</button>
      //   <input
      //     type="range"
      //     name="points"
      //     min={0}
      //     value={index}
      //     max={cluster.videos.length - 1}
      //     onChange={e => index = e.target.value }
      //   />
      //   <strong>video #{Number(index) + 1}/{cluster.videos.length} ({cluster.videos[index].duration.format('hh:mm:ss')})</strong>
      //   <h4>{cluster.videos[index].title}</h4>
      //   <published>{cluster.videos[index].publishedAt.format('dddd, MMMM Do YYYY @ h a')}</published>
      //   <span>{cluster.videos[index].description}</span>
      //
      //   <YouTube
      //     if={cluster}
      //     videoId={cluster.videos[index].youtube_id}
      //   />
      // </div>
    // )
  }


  // $ = {
  //   marginTop: 25
  // }
  //
  // $published = {
  //   fontSize: 13,
  //   marginBottom: 15,
  //   marginTop: -10,
  // }
  //
  // $span = {
  //   color: 'gray',
  //   fontSize: 13,
  //   marginTop: -10,
  //   display: 'block'
  // }
}

class ClusterThumbnail extends React.Component {

  totalDuration() {
    this.props.videos
    .map(video => video.duration)
    .reduce((total, d) => moment.duration(total).add(d), moment.duration())
  }

  render() {
    const videos = this.props.videos
    return <strong>{'cluster thumbnail'}</strong>
  }

  // $ = {
  //   // outline: '1px solid black',
  //   width: 250,
  //   fontSize: 12,
  //   margin: 25,
  //   flex: 1,
  // }
  //
  // $StackedImages = {
  //   width: 120,
  //   marginBottom: 30,
  // }
  //
  // $input = {
  //   marginLeft: -20,
  // }
  //
  // $title = {
  //   fontWeight: 'bold',
  //   fontSize: 10,
  // }
  //
  // $location = {
  //   fontSize: 12,
  // }
  //
  // $count = {
  //   display: 'inline'
  // }
  //
  // $time = {
  //   fontWeight: 500,
  //   display: 'inline',
  // }
  //
  // $stats = {
  //   textAlign: 'right',
  //   float: 'right',
  //   marginTop: -32,
  // }
}

export { ClusterExplorer, ClusterThumbnail}
