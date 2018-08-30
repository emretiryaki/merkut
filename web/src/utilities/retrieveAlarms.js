import API from 'api';

export default class AlarmList extends React.Component {
  handleSubmit = event => {
    event.preventDefault();

    API.get(`alerts`)
      .then(res => {
        console.log(res);
        console.log(res.data);
      })
  }
}