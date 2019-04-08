import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OverviewModule } from '../../overview.module';
import { FlexlayoutComponent } from './flexlayout.component';

describe('FlexlayoutComponent', () => {
  let component: FlexlayoutComponent;
  let fixture: ComponentFixture<FlexlayoutComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [OverviewModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FlexlayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});