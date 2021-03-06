// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, HostListener, ElementRef, OnInit } from '@angular/core';
import {
  Filter,
  LabelFilterService,
} from '../../services/label-filter/label-filter.service';

@Component({
  selector: 'app-input-filter',
  templateUrl: './input-filter.component.html',
  styleUrls: ['./input-filter.component.scss'],
})
export class InputFilterComponent implements OnInit {
  inputValue = '';
  showTagList = false;
  filters: Filter[] = [];

  constructor(
    private eRef: ElementRef,
    private labelFilterService: LabelFilterService
  ) {}

  ngOnInit() {
    this.labelFilterService.filters.subscribe(filters => {
      this.filters = filters;
    });
  }

  @HostListener('document:click', ['$event'])
  outsideClick(event) {
    if (!this.eRef.nativeElement.contains(event.target)) {
      this.showTagList = false;
    }
  }

  toggleTagList() {
    this.showTagList = !this.showTagList;
  }

  identifyFilter(index: number, item: Filter): string {
    return `${item.key}-${item.value}`;
  }

  remove(filter: Filter) {
    this.labelFilterService.remove(filter);
  }

  get placeholderText(): string {
    const len = this.filters.length;
    if (len > 0) {
      return `Filter by labels (${len} applied)`;
    } else {
      return 'Filter by labels';
    }
  }

  onEnter() {
    const filter = this.labelFilterService.decodeFilter(this.inputValue);
    if (filter) {
      this.labelFilterService.add(filter);
      this.inputValue = '';
      this.showTagList = true;
    } else {
      // TODO: user input value not a valid filter;
    }
  }

  clearAllFilters() {
    this.labelFilterService.clearAll();
  }
}
