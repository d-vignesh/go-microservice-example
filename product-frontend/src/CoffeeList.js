import React from 'react';
import Table from 'react-bootstrap/Table'
import axios from 'axios'
import DropdownButton from 'react-bootstrap/DropdownButton'
import Dropdown from 'react-bootstrap/Dropdown'

class CoffeeList extends React.Component {

    readData() {
        const self = this;
        axios.get(window.global.api_location+`${this.state.currencyType !== 'EUR' ?`/products?currency=${this.state.currencyType}`:'/products'}`).then(function(response) {
            console.log(response.data);
            self.setState({products: response.data});
        }).catch(function (error){
            console.log(error);
        });
    }

    getProducts() {
        let table = []

        for (let i=0; i < this.state.products.length; i++) {
            table.push(
                <tr key={i}>
                    <td>{this.state.products[i].name}</td>
                    <td>{this.state.products[i].price}</td>
                    <td>{this.state.products[i].sku}</td>
                </tr>
            );
        }

        return table
    }
    componentDidMount() {
        this.readData();
    }
    constructor(props) {
        super(props);
        this.state = {products: [], currencyType: "EUR"};

        this.readData = this.readData.bind(this)
    }

    handleOnChange = (event) => {
        console.log(event.target.value);
        this.setState({...this.state,currencyType:event.target.value}, () => this.readData())
    }

    render() {
        return (
            <div>
                <select value={this.state.currencyType} name="currencyType" id="currencyType" onChange={this.handleOnChange}>
                    <option value="USD">USD</option>
                    <option value="JPY">JPY</option>
                </select>

                <h1 style={{marginBottom: "40px"}}>Menu</h1>
                <Table>
                    <thead>
                        <tr>
                            <th>
                                Name
                            </th>
                            <th>
                                Price 
                            </th>
                            <th>
                                SKU
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {this.getProducts()}
                    </tbody>
                </Table>
            </div>
        )
    }
}

export default CoffeeList;