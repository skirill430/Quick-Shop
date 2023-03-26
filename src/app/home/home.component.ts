import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AppComponent } from '../app.component'
import { HttpService } from '../http.service';
import * as http from '../http.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent {
  title = 'Home Page';
  items : http.ItemList;
  search : string;
  showFiller = false;

  constructor(private router: Router, private httpService: HttpService) {}

  save() {
    console.log(this.search);
    this.httpService.getAllItems(this.search).subscribe(response => {
      console.log(response);
      this.items = response;
    })
  }

  // routing function to take user to pageName
  goToPage(pageName:string):void {
    this.router.navigate([`${pageName}`]);
  }
}