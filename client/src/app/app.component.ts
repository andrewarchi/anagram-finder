import { Component } from '@angular/core';
import { Http } from '@angular/http';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  providers: []
})
export class AppComponent {
  _letters: string;
  anagrams: string[];

  constructor(private http: Http) {}

  get letters(): string {
    return this._letters;
  }
  set letters(letters: string) {
    this.http.get('http://localhost:3141/anagrams/' + letters)
      .subscribe(res => this.anagrams = res.json());
    this._letters = letters;
  }
}
