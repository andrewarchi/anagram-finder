import { AnagramFinderPage } from './app.po';

describe('anagram-finder App', () => {
  let page: AnagramFinderPage;

  beforeEach(() => {
    page = new AnagramFinderPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!!');
  });
});
